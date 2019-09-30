package handler

import (
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/api"
	"github.com/micro/go-micro/api/handler"
	"github.com/micro/go-micro/api/handler/event"
	"github.com/micro/go-micro/api/router"
	"github.com/micro/go-micro/errors"

	aapi "github.com/micro/go-micro/api/handler/api"
	ahttp "github.com/micro/go-micro/api/handler/http"
	arpc "github.com/micro/go-micro/api/handler/rpc"
	aweb "github.com/micro/go-micro/api/handler/web"
)

var (
	proxyRe   = regexp.MustCompile("^[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$")
	versionRe = regexp.MustCompilePOSIX("^v[0-9]+$")
)

type metaHandler struct {
	s micro.Service
	r router.Router
}

func (m *metaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service, err := m.r.Route(r)
	if err != nil {
		er := errors.InternalServerError(m.r.Options().Namespace, err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(er.Error()))
		return
	}

	//use for log
	r.Header["Request"] = []string{r.Method, r.RequestURI, r.Proto}

	//check and update handler type
	m.checkHandlerType(service)

	//dispatch
	switch service.Endpoint.Handler {
	case aapi.Handler:
		aapi.WithService(service, handler.WithService(m.s)).ServeHTTP(w, r)
	case aweb.Handler:
		aweb.WithService(service, handler.WithService(m.s)).ServeHTTP(w, r)
	case "proxy", ahttp.Handler:
		ahttp.WithService(service, handler.WithService(m.s)).ServeHTTP(w, r)
	case arpc.Handler:
		arpc.WithService(service, handler.WithService(m.s)).ServeHTTP(w, r)
	case event.Handler:
		ev := event.NewHandler(
			handler.WithNamespace(m.r.Options().Namespace),
			handler.WithService(m.s),
		)
		ev.ServeHTTP(w, r)
	default:
		arpc.WithService(service, handler.WithService(m.s)).ServeHTTP(w, r)
	}
}

//check handler type use metadata
//add handler type to metadata when registry
func (m *metaHandler) checkHandlerType(s *api.Service) {
	if s == nil {
		return
	}

	if s.Endpoint == nil {
		return
	}

	if len(s.Services) == 0 {
		return
	}

	service := s.Services[0]
	for _, v := range service.Endpoints {
		if v.Name != s.Endpoint.Name {
			//convert
			name := m.handlerRoute(v.Name)
			if name != s.Endpoint.Name {
				continue
			}
		}

		handler, ok := v.Metadata["handler"]
		if !ok {
			return
		}

		s.Endpoint.Handler = handler
		return
	}
}


//convert method，eg:/echo/hi -> Echo.Hi
func (m *metaHandler) handlerRoute(p string) string {
	p = path.Clean(p)
	p = strings.TrimPrefix(p, "/")
	parts := strings.Split(p, "/")

	if len(parts) <= 2 {
		return methodName(parts)
	}

	return methodName(parts[len(parts)-2:])
}

func methodName(parts []string) string {
	for i, part := range parts {
		parts[i] = toCamel(part)
	}

	return strings.Join(parts, ".")
}

func toCamel(s string) string {
	words := strings.Split(s, "-")
	var out string
	for _, word := range words {
		out += strings.Title(word)
	}
	return out
}


// Meta is a http.Handler that routes based on endpoint metadata
func Meta(s micro.Service, r router.Router) http.Handler {
	return &metaHandler{ s: s, r: r }
}
