package html_test

import (
	"bytes"
	"net/http"
	"os"
	"regexp"
	"strings"

	. "github.com/BeamFoundry/rays-fiber/pkg/template/html"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	admin         = "admin"
	complexexpect = `<!DOCTYPE html><html><head><title>Main</title></head><body><h2>Header</h2><h1>Hello, World!</h1><h2>Footer</h2></body></html>`
)

func trim(str string) string {
	trimmed := strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(str, " "))
	trimmed = strings.ReplaceAll(trimmed, " <", "<")
	trimmed = strings.ReplaceAll(trimmed, "> ", ">")
	return trimmed
}

var _ = Describe("HTML Template", func() {
	It("can Render", func() {
		engine := New("./views", ".html")
		engine.AddFunc("isAdmin", func(user string) bool {
			return user == admin
		})
		Expect(engine.Load()).NotTo(HaveOccurred())

		// Partials
		var buf bytes.Buffer
		err := engine.Render(&buf, "index", map[string]any{
			"Title": "Hello, World!",
		}, "partials/header", "partials/footer")
		Expect(err).NotTo(HaveOccurred())

		expect := `<h2>Header</h2><h1>Hello, World!</h1><h2>Footer</h2>`
		result := trim(buf.String())
		Expect(expect).To(Equal(result))

		// Single
		buf.Reset()
		err = engine.Render(&buf, "errors/404", map[string]any{
			"Error": "404 Not Found!",
		})
		Expect(err).NotTo(HaveOccurred())

		expect = `<h1>404 Not Found!</h1>`
		result = trim(buf.String())
		Expect(expect).To(Equal(result))
	})

	It("can AddFunc", func() {
		engine := New("./views", ".html")
		engine.AddFunc("isAdmin", func(user string) bool {
			return user == admin
		})
		Expect(engine.Load()).NotTo(HaveOccurred())

		// Func is admin
		var buf bytes.Buffer
		err := engine.Render(&buf, admin, map[string]any{
			"User": admin,
		})
		Expect(err).NotTo(HaveOccurred())

		expect := `<h1>Hello, Admin!</h1>`
		result := trim(buf.String())
		Expect(expect).To(Equal(result))

		// Func is not admin
		buf.Reset()
		err = engine.Render(&buf, admin, map[string]any{
			"User": "john",
		})
		Expect(err).NotTo(HaveOccurred())

		expect = `<h1>Access denied!</h1>`
		result = trim(buf.String())
		Expect(expect).To(Equal(result))

		// FuncMap
		fm := engine.FuncMap()
		_, ok := fm["isAdmin"]
		Expect(ok).To(BeTrue())
	})

	It("can AddFuncMap", func() {
		// Create a temporary directory
		dir, err := os.MkdirTemp(".", "")
		Expect(err).NotTo(HaveOccurred())

		defer func() {
			Expect(
				os.RemoveAll(dir),
			).NotTo(HaveOccurred())
		}()

		// Create a temporary template file.
		err = os.WriteFile(dir+"/func_map.html", []byte(`<h2>{{lower .Var1}}</h2><p>{{upper .Var2}}</p>`), 0o600)
		Expect(err).NotTo(HaveOccurred())

		engine := New(dir, ".html")
		fm := map[string]any{
			"lower": strings.ToLower,
			"upper": strings.ToUpper,
		}

		engine.AddFuncMap(fm)
		Expect(engine.Load()).NotTo(HaveOccurred())

		var buf bytes.Buffer
		err = engine.Render(&buf, "func_map", map[string]any{
			"Var1": "LOwEr",
			"Var2": "upPEr",
		})
		Expect(err).NotTo(HaveOccurred())

		expect := `<h2>lower</h2><p>UPPER</p>`
		result := trim(buf.String())
		Expect(expect).To(Equal(result))

		// FuncMap
		fm2 := engine.FuncMap()
		_, ok := fm2["lower"]
		Expect(ok).To(BeTrue())
		_, ok = fm2["upper"]
		Expect(ok).To(BeTrue())
	})

	It("can Layout", func() {
		engine := New("./views", ".html")

		engine.AddFunc("isAdmin", func(user string) bool {
			return user == admin
		})
		Expect(engine.Load()).NotTo(HaveOccurred())

		var buf bytes.Buffer
		err := engine.Render(&buf, "index", map[string]any{
			"Title":    "Hello, World!",
			"Template": "layouts/main",
		}, "layouts/main", "partials/header", "partials/footer")
		Expect(err).NotTo(HaveOccurred())

		result := trim(buf.String())
		Expect(complexexpect).To(Equal(result))
	})

	It("can Empty_Layout", func() {
		engine := New("./views", ".html")

		engine.AddFunc("isAdmin", func(user string) bool {
			return user == admin
		})
		Expect(engine.Load()).NotTo(HaveOccurred())

		var buf bytes.Buffer
		err := engine.Render(&buf, "index", map[string]any{
			"Title": "Hello, World!",
		}, "partials/header", "partials/footer", "")
		Expect(err).NotTo(HaveOccurred())

		expect := `<h2>Header</h2><h1>Hello, World!</h1><h2>Footer</h2>`
		result := trim(buf.String())
		Expect(expect).To(Equal(result))
	})

	// Test_Layout_Multi checks if the LayoutName can be rendered multiple times
	It("can Layout_Multi", func() {
		engine := New("./views", ".html")

		engine.AddFunc("isAdmin", func(user string) bool {
			return user == admin
		})
		Expect(engine.Load()).NotTo(HaveOccurred())

		for range 2 {
			var buf bytes.Buffer
			err := engine.Render(&buf, "index", map[string]any{
				"Title":    "Hello, World!",
				"Template": "layouts/main",
			}, "layouts/main", "partials/header", "partials/footer")
			Expect(err).NotTo(HaveOccurred())

			result := trim(buf.String())
			Expect(complexexpect).To(Equal(result))
		}
	})

	It("can Layout_Nested", func() {
		engine := New("./views", ".html")

		engine.AddFunc("isAdmin", func(user string) bool {
			return user == admin
		})
		Expect(engine.Load()).NotTo(HaveOccurred())

		var buf bytes.Buffer
		err := engine.Render(&buf, "index", map[string]any{
			"Title":    "Hello, World!",
			"Template": "layouts/nested/main",
		}, "layouts/nested/main", "layouts/nested/base", "partials/header", "partials/footer")
		Expect(err).NotTo(HaveOccurred())

		result := trim(buf.String())
		Expect(complexexpect).To(Equal(result))
	})

	It("can FileSystem", func() {
		engine := NewFileSystem(http.Dir("./views"), ".html")

		engine.AddFunc("isAdmin", func(user string) bool {
			return user == admin
		})
		Expect(engine.Load()).NotTo(HaveOccurred())

		var buf bytes.Buffer
		err := engine.Render(&buf, "index", map[string]any{
			"Title":    "Hello, World!",
			"Template": "layouts/main",
		}, "layouts/main", "partials/header", "partials/footer")
		Expect(err).NotTo(HaveOccurred())

		result := trim(buf.String())
		Expect(complexexpect).To(Equal(result))
	})

	It("can Reload", func() {
		engine := NewFileSystem(http.Dir("./views"), ".html")
		engine.Reload(true) // Optional. Default: false

		engine.AddFunc("isAdmin", func(user string) bool {
			return user == admin
		})
		Expect(engine.Load()).NotTo(HaveOccurred())

		Expect(
			os.WriteFile("./views/ShouldReload.html", []byte("after ShouldReload\n"), 0o600),
		).NotTo(HaveOccurred())
		defer func() {
			Expect(
				os.WriteFile("./views/ShouldReload.html", []byte("before ShouldReload\n"), 0o600),
			).NotTo(HaveOccurred())
		}()

		Expect(engine.Load()).NotTo(HaveOccurred())

		var buf bytes.Buffer
		Expect(
			engine.Render(&buf, "ShouldReload", nil),
		).NotTo(HaveOccurred())

		expect := "after ShouldReload"
		result := trim(buf.String())
		Expect(expect).To(Equal(result))
	})
})

/*
func Benchmark_Html(b *testing.B) {
	expectSimple := `<h1>Hello, World!</h1>`
	expectExtended := `<!DOCTYPE html><html><head><title>Main</title></head><body><h2>Header</h2><h1>Hello, Admin!</h1><h2>Footer</h2></body></html>`
	engine := New("./views", ".html")
	engine.AddFunc("isAdmin", func(user string) bool {
		return user == admin
	})
	require.NoError(b, engine.Load())

	b.Run("simple", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		var buf bytes.Buffer
		for i := 0; i < bb.N; i++ {
			buf.Reset()
			//nolint:gosec,errcheck // Return value not needed for benchmark
			_ = engine.Render(&buf, "simple", map[string]any{
				"Title": "Hello, World!",
			})
		}
	})

	b.Run("extended", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		var buf bytes.Buffer
		for i := 0; i < bb.N; i++ {
			buf.Reset()
			//nolint:gosec,errcheck // Return value not needed for benchmark
			_ = engine.Render(&buf, "extended", map[string]any{
				"User": admin,
			}, "layouts/main")
		}
	})

	b.Run("simple_asserted", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		var buf bytes.Buffer
		for i := 0; i < bb.N; i++ {
			buf.Reset()
			err := engine.Render(&buf, "simple", map[string]any{
				"Title": "Hello, World!",
			})

			require.NoError(b, err)
			require.Equal(b, expectSimple, trim(buf.String()))
		}
	})

	b.Run("extended_asserted", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		var buf bytes.Buffer
		for i := 0; i < bb.N; i++ {
			buf.Reset()
			err := engine.Render(&buf, "extended", map[string]any{
				"User": admin,
			}, "layouts/main")

			require.NoError(b, err)
			require.Equal(b, expectExtended, trim(buf.String()))
		}
	})
})

func Benchmark_Html_Parallel(b *testing.B) {
	expectSimple := `<h1>Hello, Parallel!</h1>`
	expectExtended := `<!DOCTYPE html><html><head><title>Main</title></head><body><h2>Header</h2><h1>Hello, Admin!</h1><h2>Footer</h2></body></html>`
	engine := New("./views", ".html")
	engine.AddFunc("isAdmin", func(user string) bool {
		return user == admin
	})
	require.NoError(b, engine.Load())

	b.Run("simple", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				var buf bytes.Buffer
				//nolint:gosec,errcheck // Return value not needed for benchmark
				_ = engine.Render(&buf, "simple", map[string]any{
					"Title": "Hello, Parallel!",
				})
			}
		})
	})

	b.Run("extended", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				var buf bytes.Buffer
				//nolint:gosec,errcheck // Return value not needed for benchmark
				_ = engine.Render(&buf, "extended", map[string]any{
					"User": admin,
				}, "layouts/main")
			}
		})
	})

	b.Run("simple_asserted", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				var buf bytes.Buffer
				err := engine.Render(&buf, "simple", map[string]any{
					"Title": "Hello, Parallel!",
				})
				require.NoError(bb, err)
				require.Equal(bb, expectSimple, trim(buf.String()))
			}
		})
	})

	b.Run("extended_asserted", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				var buf bytes.Buffer
				err := engine.Render(&buf, "extended", map[string]any{
					"User": admin,
				}, "layouts/main")
				require.NoError(bb, err)
				require.Equal(bb, expectExtended, trim(buf.String()))
			}
		})
	})
})
*/
