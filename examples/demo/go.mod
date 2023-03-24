module github.com/MuggleWei/goev/examples/demo

go 1.19

replace github.com/MuggleWei/goev => ../../

require (
	github.com/MuggleWei/goev v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
