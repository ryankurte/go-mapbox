package maps

type MapID string

const (
	MapIDStreets          MapID = "mapbox.streets"
	MapIDLight            MapID = "mapbox.light"
	MapIDDark             MapID = "mapbox.dark"
	MapIDSatellite        MapID = "mapbox.satellite"
	MapIDStreetsSatellite MapID = "mapbox.streets-satellite"
	MapIDWheatpaste       MapID = "mapbox.wheatpaste"
	MapIDStreetsBasic     MapID = "mapbox.streets-basic"
	MapIDComic            MapID = "mapbox.comic"
	MapIDOutdoors         MapID = "mapbox.outdoors"
	MapIDRunBikeHike      MapID = "mapbox.run-bike-hike"
	MapIDPencil           MapID = "mapbox.pencil"
	MapIDPirates          MapID = "mapbox.pirates"
	MapIDEmerald          MapID = "mapbox.emerald"
	MapIDHighContrast     MapID = "mapbox.high-contrast"
	MapIDTerrainRGB       MapID = "mapbox.terrain-rgb"
)

type MapFormat string

const (
	MapFormatPng    MapFormat = "png"    // true color PNG
	MapFormatPng32  MapFormat = "png32"  // 32 color indexed PNG
	MapFormatPng64  MapFormat = "png64"  // 64 color indexed PNG
	MapFormatPng128 MapFormat = "png128" // 128 color indexed PNG
	MapFormatPng256 MapFormat = "png256" // 256 color indexed PNG
	MapFormatPngRaw MapFormat = "pngraw" // Raw PNG (only for MapIDTerrainRGB)
	MapFormatJpg70  MapFormat = "jpg70"  // 70% quality JPG
	MapFormatJpg80  MapFormat = "jpg80"  // 80% quality JPG
	MapFormatJpg90  MapFormat = "jpg90"  // 90% quality JPG
)
