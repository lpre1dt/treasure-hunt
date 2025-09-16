package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jomei/notionapi"
)

//go:embed templates/*
var templates embed.FS

// App enthält alle App-Komponenten
type App struct {
	notion        *notionapi.Client
	teamsDBID     string
	challengeDBID string
	templates     *template.Template
}

func main() {
	// .env Datei laden
	if err := godotenv.Load(); err != nil {
		log.Println("Keine .env Datei gefunden, nutze Umgebungsvariablen")
	}

	// Konfiguration aus Umgebung
	notionToken := os.Getenv("NOTION_TOKEN")
	teamsDBID := os.Getenv("TEAMS_DB_ID")
	challengeDBID := os.Getenv("CHALLENGES_DB_ID")

	if notionToken == "" || teamsDBID == "" || challengeDBID == "" {
		log.Fatal("NOTION_TOKEN, TEAMS_DB_ID und CHALLENGES_DB_ID müssen in .env gesetzt sein")
	}

	// Notion Client initialisieren
	client := notionapi.NewClient(notionapi.Token(notionToken))

	// Templates laden
	tmpl, err := template.ParseFS(templates, "templates/*.html")
	if err != nil {
		log.Fatal("Fehler beim Laden der Templates:", err)
	}

	app := &App{
		notion:        client,
		teamsDBID:     teamsDBID,
		challengeDBID: challengeDBID,
		templates:     tmpl,
	}

	// Gin Router einrichten
	r := gin.Default()

	// Routes
	r.GET("/", app.handleHome)
	r.GET("/next/:id", app.handleChallengeForm)
	r.POST("/next/:id", app.handleNextChallenge)
	r.GET("/mvpgenerator", app.handleMVPGenerator)

	// Server starten
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server startet auf http://localhost:%s", port)
	r.Run(":" + port)
}

// handleHome zeigt Startseite
func (app *App) handleHome(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := app.templates.ExecuteTemplate(c.Writer, "home.html", nil); err != nil {
		c.String(http.StatusInternalServerError, "Template-Fehler: %v", err)
	}
}

// getAllTeamNames holt alle Teamnamen aus der Notion DB
func (app *App) getAllTeamNames() ([]string, error) {
	ctx := context.Background()
	var teamNames []string

	// Query, um alle Seiten aus der Team-DB zu holen
	query := &notionapi.DatabaseQueryRequest{
		PageSize: 100, // Annahme: Es gibt nicht mehr als 100 Teams
	}

	result, err := app.notion.Database.Query(ctx, notionapi.DatabaseID(app.teamsDBID), query)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Abfragen der Team-Datenbank: %w", err)
	}

	// Iteriere durch die Ergebnisse und extrahiere den Titel jeder Seite
	for _, page := range result.Results {
		// Die Titel-Eigenschaft hat keinen festen Namen, sie wird durch ihren Typ identifiziert.
		for _, prop := range page.Properties {
			if titleProp, ok := prop.(*notionapi.TitleProperty); ok {
				if len(titleProp.Title) > 0 {
					teamNames = append(teamNames, titleProp.Title[0].PlainText)
					break // Nächste Seite, da wir den Titel gefunden haben
				}
			}
		}
	}

	if len(teamNames) == 0 {
		return nil, fmt.Errorf("keine Teams in der Datenbank gefunden")
	}

	// Sortiere die Teamnamen alphabetisch
	sort.Strings(teamNames)

	return teamNames, nil
}
func (app *App) handleMVPGenerator(c *gin.Context) {

	// MVP-Generierung starten
	app.generateMVP(c)

}
func (app *App) generateMVP(c *gin.Context) {
	// Beispiel-MVP-Generierung (hier könnte eine komplexere Logik stehen)
	mvpArray := []string{
		"Portable Water Filter",
		"Eco-Friendly Phone Stand",
		"Leaf-Based Notebook",
		"Natural Air Freshener",
		"Solar-Powered Lantern",
		"Pocket-Sized Board Game",
		"Mini Bird Feeder",
		"Portable Hammock",
		"Eco Speaker Amplifier",
		"Handmade Jewelry",
		"Biodegradable Straw",
		"Reusable Cutlery Set",
		"Eco Toy Car",
		"Compostable Food Container",
		"Emergency Shelter Kit",
		"Natural Bandage",
		"Seed-Paper Business Card",
		"Upcycled Coin Wallet",
		"Outdoor Chess Set",
		"DIY Kite",
		"Portable Plant Pot",
		"Zero-Waste Picnic Kit",
		"Wind Chime",
		"Stress Relief Toy",
		"Eco Bracelet",
		"Paper Recycling Kit",
		"Biodegradable Soap Holder",
		"Natural Toothbrush",
		"Eco-Friendly Bag",
		"Outdoor Survival Kit",
		"Pocket Garden",
		"Park-Themed Board Game",
		"Compost Bin Prototype",
		"Mini Solar Oven",
		"Eco-Friendly Candle",
		"Toy Drone Shell",
		"Bird Call Whistle",
		"Outdoor Gym Equipment",
		"Eco Keychain",
		"Rainwater Collector",
		"DIY Frisbee",
		"Eco Sunglasses",
		"Reusable Coffee Sleeve",
		"Portable Charger Holder",
		"Nature Bookmark",
		"Toy Boat",
		"DIY Musical Instrument",
		"Eco-Friendly Packaging",
		"Foldable Stool",
		"Eco Travel Mug",
		"Mini Wind Turbine",
		"Outdoor Meditation Mat",
		"Portable Whiteboard",
		"Emergency Cooking Stove",
		"Toy Puzzle",
		"Nature Art Frame",
		"Eco Phone Case",
		"Reusable Water Filter Straw",
		"Picnic Blanket Prototype",
		"Eco Bag Tag",
		"DIY Pen Holder",
		"Eco-Friendly Wallet",
		"Upcycled Backpack",
		"Biodegradable Plant Pot",
		"Eco Toothpaste Holder",
		"DIY Notebook",
		"Park Cleaning Kit",
		"Eco Coaster",
		"Portable Light Reflector",
		"Outdoor Cooking Kit",
		"Eco-Friendly Umbrella",
		"Eco-Friendly Badge",
		"Mini Greenhouse",
		"Eco-Friendly Soap Dish",
		"Outdoor Relaxation Chair",
		"Toy Rocket",
		"Eco-Friendly Speaker",
		"DIY Lamp",
		"Eco-Friendly Calendar",
		"Nature Camera Case",
		"Eco-Friendly Shoes",
		"Eco-Friendly Gloves",
		"Outdoor Game Dice",
		"DIY Jewelry Box",
		"Eco Candle Holder",
		"Portable Fire Starter",
		"Eco Ashtray",
		"Reusable Straw Holder",
		"DIY Sunglass Holder",
		"Eco Plant Sprayer",
		"Toy Airplane",
		"Mini Compost Bag",
		"Eco-Friendly Watch",
		"Portable Raincoat",
		"Eco Pencil Case",
		"Outdoor Card Game",
		"Toy Binoculars",
		"Eco Lantern",
		"Pocket First Aid Kit",
		"Eco Water Bottle",
		"Eco Blanket",
	}
	// Zufälliges MVP auswählen
	mvp := mvpArray[randRange(0, len(mvpArray))]

	// Ergebnis an das Template übergeben
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := app.templates.ExecuteTemplate(c.Writer, "mvpgenerator.html", gin.H{
		"mvp": mvp,
	}); err != nil {
		c.String(http.StatusInternalServerError, "Template-Fehler: %v", err)
	}
}
func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}

// handleChallengeForm zeigt Formular für Teamname-Eingabe mit Dropdown
func (app *App) handleChallengeForm(c *gin.Context) {
	challengeID := c.Param("id")

	// Alle Teamnamen aus Notion für das Dropdown holen
	teamNames, err := app.getAllTeamNames()
	if err != nil {
		log.Printf("Fehler beim Abrufen der Teamnamen: %v", err)
		c.String(http.StatusInternalServerError, "Fehler beim Laden der Teamliste: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := app.templates.ExecuteTemplate(c.Writer, "teamform.html", gin.H{
		"challengeID": challengeID,
		"Teams":       teamNames, // Teamliste an das Template übergeben
	}); err != nil {
		c.String(http.StatusInternalServerError, "Template-Fehler: %v", err)
	}
}

// handleNextChallenge verarbeitet Team und leitet zur nächsten Challenge weiter
func (app *App) handleNextChallenge(c *gin.Context) {
	currentChallengeID := c.Param("id")
	teamName := c.PostForm("team")

	if teamName == "" {
		c.String(http.StatusBadRequest, "Teamname erforderlich")
		return
	}

	log.Printf("Suche Team: %s mit Challenge ID: %s", teamName, currentChallengeID)

	// Finde Team-Page in Teams-DB
	teamPageID, err := app.findTeamPage(teamName)
	if err != nil || teamPageID == "" {
		log.Printf("Team nicht gefunden: %s", teamName)
		c.Header("Content-Type", "text/html; charset=utf-8")
		app.templates.ExecuteTemplate(c.Writer, "error.html", gin.H{
			"error": "Team nicht gefunden",
		})
		return
	}

	log.Printf("Team-Page gefunden: %s", teamPageID)

	// Hole Team-Daten
	teamData, err := app.getTeamChallenges(teamPageID)
	if err != nil {
		log.Printf("Fehler beim Abrufen der Challenges: %v", err)
		c.Header("Content-Type", "text/html; charset=utf-8")
		app.templates.ExecuteTemplate(c.Writer, "error.html", gin.H{
			"error": fmt.Sprintf("Fehler beim Abrufen der Team-Daten: %v", err),
		})
		return
	}

	log.Printf("Gefundene Challenges: %v", teamData)

	// Finde aktuelle Challenge-Position und hole nächste
	nextChallengeURL := app.findNextChallengeURL(teamData, currentChallengeID)

	if nextChallengeURL == "" {
		log.Printf("Keine weitere Challenge gefunden nach ID: %s", currentChallengeID)
		// Keine weitere Challenge oder Challenge nicht gefunden
		c.Header("Content-Type", "text/html; charset=utf-8")
		app.templates.ExecuteTemplate(c.Writer, "finished.html", gin.H{
			"team": teamName,
		})
		return
	}

	log.Printf("Nächste Challenge URL: %s", nextChallengeURL)

	// Weiterleitung zur nächsten Challenge
	c.Header("Content-Type", "text/html; charset=utf-8")
	app.templates.ExecuteTemplate(c.Writer, "redirect.html", gin.H{
		"url":  nextChallengeURL,
		"team": teamName,
	})
}

// findTeamPage findet die Team-Page ID anhand des Teamnamens
func (app *App) findTeamPage(teamName string) (string, error) {
	ctx := context.Background()

	// Versuche verschiedene Property-Namen für den Titel
	propertyNames := []string{"Name", "Team", "Title", "title"}

	for _, prop := range propertyNames {
		filter := &notionapi.DatabaseQueryRequest{
			Filter: &notionapi.PropertyFilter{
				Property: prop,
				RichText: &notionapi.TextFilterCondition{
					Equals: teamName,
				},
			},
		}

		result, err := app.notion.Database.Query(ctx, notionapi.DatabaseID(app.teamsDBID), filter)
		if err == nil && len(result.Results) > 0 {
			return string(result.Results[0].ID), nil
		}
	}

	// Als letzten Versuch: Ohne Filter alle Teams holen und manuell suchen
	allTeams := &notionapi.DatabaseQueryRequest{
		PageSize: 100,
	}

	result, err := app.notion.Database.Query(ctx, notionapi.DatabaseID(app.teamsDBID), allTeams)
	if err != nil {
		return "", err
	}

	// Manuell nach dem Team suchen
	for _, page := range result.Results {
		for _, prop := range page.Properties {
			switch p := prop.(type) {
			case *notionapi.TitleProperty:
				if len(p.Title) > 0 && p.Title[0].PlainText == teamName {
					return string(page.ID), nil
				}
			case *notionapi.RichTextProperty:
				if len(p.RichText) > 0 && p.RichText[0].PlainText == teamName {
					return string(page.ID), nil
				}
			}
		}
	}

	return "", nil
}

// getTeamChallenges holt alle Challenge-Relations einer Team-Page
func (app *App) getTeamChallenges(teamPageID string) (map[int]string, error) {
	ctx := context.Background()

	page, err := app.notion.Page.Get(ctx, notionapi.PageID(teamPageID))
	if err != nil {
		return nil, err
	}

	challenges := make(map[int]string)

	// Durchsuche alle Properties nach Challenge-Relations
	for propName, prop := range page.Properties {
		// Prüfe ob Property eine Challenge-Relation ist (Challenge1, Challenge2, etc.)
		var challengeNum int
		if _, err := fmt.Sscanf(propName, "Challenge%d", &challengeNum); err == nil {
			if relation, ok := prop.(*notionapi.RelationProperty); ok && len(relation.Relation) > 0 {
				// Hole die verlinkte Challenge-Page
				challengePageID := string(relation.Relation[0].ID)
				challengePage, err := app.notion.Page.Get(ctx, notionapi.PageID(challengePageID))
				if err == nil {
					// Extrahiere Challenge-ID aus der "id" Property (Number)
					if idProp, ok := challengePage.Properties["id"]; ok {
						switch v := idProp.(type) {
						case *notionapi.NumberProperty:
							challenges[challengeNum] = fmt.Sprintf("%.0f", v.Number)
						}
					}
					// Fallback: Versuche "ID" (groß)
					if _, exists := challenges[challengeNum]; !exists {
						if idProp, ok := challengePage.Properties["ID"]; ok {
							switch v := idProp.(type) {
							case *notionapi.NumberProperty:
								challenges[challengeNum] = fmt.Sprintf("%.0f", v.Number)
							}
						}
					}
				}
			}
		}
	}

	return challenges, nil
}

// findNextChallengeURL findet die URL der nächsten Challenge
func (app *App) findNextChallengeURL(challenges map[int]string, currentID string) string {
	// Finde Position der aktuellen Challenge
	currentPos := 0
	for pos, id := range challenges {
		if id == currentID {
			currentPos = pos
			break
		}
	}

	// Suche nächste Challenge (currentPos + 1)
	nextPos := currentPos + 1
	if nextID, exists := challenges[nextPos]; exists {
		// Hole die Challenge-Page aus der Challenge-DB mit der ID
		ctx := context.Background()

		// Suche Challenge mit dieser ID (Number Property)
		filter := &notionapi.DatabaseQueryRequest{
			Filter: &notionapi.PropertyFilter{
				Property: "id",
				Number: &notionapi.NumberFilterCondition{
					Equals: func() *float64 {
						f, _ := parseFloat(nextID)
						return &f
					}(),
				},
			},
		}

		result, err := app.notion.Database.Query(ctx, notionapi.DatabaseID(app.challengeDBID), filter)

		// Fallback: Versuche "ID" (groß)
		if err != nil || len(result.Results) == 0 {
			filter.Filter = &notionapi.PropertyFilter{
				Property: "ID",
				Number: &notionapi.NumberFilterCondition{
					Equals: func() *float64 {
						f, _ := parseFloat(nextID)
						return &f
					}(),
				},
			}
			result, _ = app.notion.Database.Query(ctx, notionapi.DatabaseID(app.challengeDBID), filter)
		}

		if result != nil && len(result.Results) > 0 {
			// Baue Notion-URL
			pageID := string(result.Results[0].ID)
			// Entferne Bindestriche aus der ID für die URL
			cleanID := ""
			for _, c := range pageID {
				if c != '-' {
					cleanID += string(c)
				}
			}
			return fmt.Sprintf("https://marcbaumholz.notion.site/%s", cleanID)
		}
	}

	return ""
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}
