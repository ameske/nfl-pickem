build:
	cd cli && go build -o nfl && mv nfl ..
	cd server && go build -o nflserver && mv nflserver ..

install: build
	mv opt/ameske/gonfl/templates /opt/ameske/gonfl/templates
	mv nfl /opt/ameske/gonfl/nfl
	mv nflserver /opt/ameske/gonfl/nflserver

uninstall:
	rm -r /opt/ameske/gonfl
	
clean:
	rm nfl
	rm nflserver
