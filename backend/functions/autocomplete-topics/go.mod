module github.com/raphink/debate/backend/functions/autocomplete-topics

go 1.24

require (
	cloud.google.com/go/firestore v1.18.0
	github.com/GoogleCloudPlatform/functions-framework-go v1.9.0
	github.com/raphink/debate/backend/shared v0.0.0
	google.golang.org/api v0.215.0
)

replace github.com/raphink/debate/backend/shared => ../../shared
