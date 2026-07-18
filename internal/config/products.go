// Package config holds the catalogue of supported JetBrains products.
package config

// Canonical JetBrains product identifiers.
const (
	IntelliJIdea = "IntelliJIdea"
	CLion        = "CLion"
	PhpStorm     = "PhpStorm"
	GoLand       = "GoLand"
	PyCharm      = "PyCharm"
	WebStorm     = "WebStorm"
	Rider        = "Rider"
	DataGrip     = "DataGrip"
	RubyMine     = "RubyMine"
	AppCode      = "AppCode"
)

// Products lists every supported JetBrains IDE.
var Products = []string{ // const
	IntelliJIdea, CLion, PhpStorm, GoLand, PyCharm,
	WebStorm, Rider, DataGrip, RubyMine, AppCode,
}

// TokenToProduct maps a JVM selector/prefix token to its canonical product.
var TokenToProduct = map[string]string{ // const
	"Idea": IntelliJIdea, "CLion": CLion, "PhpStorm": PhpStorm,
	"GoLand": GoLand, "PyCharm": PyCharm, "WebStorm": WebStorm,
	"Rider": Rider, "DataGrip": DataGrip, "RubyMine": RubyMine,
	"AppCode": AppCode,
}

// ProductExe maps each product to its Windows launcher image name.
var ProductExe = map[string]string{ // const
	IntelliJIdea: "idea64.exe", CLion: "clion64.exe", PhpStorm: "phpstorm64.exe",
	GoLand: "goland64.exe", PyCharm: "pycharm64.exe", WebStorm: "webstorm64.exe",
	Rider: "rider64.exe", DataGrip: "datagrip64.exe", RubyMine: "rubymine64.exe",
	AppCode: "appcode64.exe",
}
