cd game
go build
move game.exe ..

cd ..

cd gate
go build
move gate.exe ..

cd ..

cd portal
go build
move portal.exe ..

cd ..

cd allinone
go build
move allinone.exe ..

ping /n 10 127.0 > NUL