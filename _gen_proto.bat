cd proto

cd inner
protoc --go_out=. *.proto
go install

cd ../nfy
protoc --go_out=. *.proto
protoc -onfy.pb nfy.proto
go install
move nfy.pb ../..

cd ../rsp
protoc --go_out=. *.proto
protoc -orsp.pb rsp.proto
go install
move rsp.pb ../..

cd ../req
protoc --go_out=. *.proto
protoc -oreq.pb req.proto
go install
move req.pb ../..

cd ../..
move nfy.pb ..\..\..\..\..\Unity3D\SimpleFramework_UGUI-master\Assets\Lua\Proto
move rsp.pb ..\..\..\..\..\Unity3D\SimpleFramework_UGUI-master\Assets\Lua\Proto
move req.pb ..\..\..\..\..\Unity3D\SimpleFramework_UGUI-master\Assets\Lua\Proto

ping /n 10 127.0 > NUL