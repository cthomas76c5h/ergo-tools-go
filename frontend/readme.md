# in elm-frontend/
elm make src/Main.elm --output=public/elm.js
# serve the static files (pick one)
npx serve public
# or
python -m http.server --directory public 5173
