package swagger

import (
	"log"
	"net/http"
	"path"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

//ServeSwaggerFile 把proto文件夹中的swagger.json文件暴露出去
func ServeSwaggerFile(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "swagger.json") {
		log.Printf("Not Found: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	// "../proto/"为.swagger.json所在目录
	p = path.Join("../proto/", p)

	log.Printf("Serving swagger-file: %s", p)

	http.ServeFile(w, r, p)
}

//ServeSwaggerUI 对外提供swagger-ui
func ServeSwaggerUI(mux *http.ServeMux) {
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "swagger/swagger-ui", //swagger-ui文件夹所在目录
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}
