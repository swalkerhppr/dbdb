package assets

import "embed"

//go:embed game_assets/asset_definitions.xml
var metadata []byte

//go:embed game_assets
var assetFS embed.FS
