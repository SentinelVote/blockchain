# This patch updates the image's maximum request body size to 500MB,
# so that it can accept a foldedPublicKeys string larger than 2000 users.
FROM softwaremill/fablo-rest:0.1.0
RUN sed -i.bak 's/app.use(express_1.default.json({ type: () => "json" }))/app.use(express_1.default.json({ type: () => "json", limit: "500mb" }))/g' index.js
