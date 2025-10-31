package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// NOTE: Drink struct, drinks variable, et getMenu ont été déplacés dans drink.go

// Middleware CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func testPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Brew & Co - Commandez vos boissons</title>
		<style>
			:root {
				--coffee-dark: #4A3A2A; /* Marron Foncé */
				--background: #F8F5F0; /* Fond Crème */
				--accent: #D8B468; /* Accent Caramel/Or */
				--text-color: #333;
			}
			body {
				font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
				margin: 0;
				padding: 0;
				background-color: var(--background);
				color: var(--text-color);
			}
			.container {
				max-width: 960px;
				margin: 40px auto;
				padding: 30px;
				background-color: white;
				box-shadow: 0 10px 25px rgba(0,0,0,0.1);
				border-radius: 12px;
			}
			header {
				background-color: var(--coffee-dark);
				color: white;
				padding: 20px 0;
				margin: -30px -30px 30px -30px; /* Extend header to edges of container padding */
				border-top-left-radius: 12px;
				border-top-right-radius: 12px;
				text-align: center;
			}
			h1 {
				margin: 0;
				font-size: 2em;
			}
			section {
				margin-bottom: 30px;
				border: 1px solid #E0E0E0;
				padding: 20px;
				border-radius: 8px;
				background-color: #FFFFFF;
				box-shadow: 0 2px 5px rgba(0,0,0,0.05);
			}
			h2 {
				color: var(--coffee-dark);
				border-bottom: 2px solid var(--accent);
				padding-bottom: 5px;
				margin-top: 0;
				font-size: 1.5em;
			}
			button, input[type="submit"] {
				padding: 10px 20px;
				background: var(--accent);
				color: var(--coffee-dark);
				border: none;
				border-radius: 5px;
				cursor: pointer;
				transition: background 0.3s, transform 0.1s;
				font-size: 16px;
				font-weight: bold;
			}
			button:hover, input[type="submit"]:hover {
				background: #C3A158; /* Slightly darker accent */
				transform: translateY(-1px);
			}
			pre {
				background: #333;
				color: #4CAF50; /* Texte vert console */
				padding: 15px;
				border-radius: 5px;
				overflow-x: auto;
				white-space: pre-wrap;
				word-break: break-all;
				font-size: 0.9em;
			}
			form {
				display: flex;
				flex-direction: column;
				gap: 15px;
			}
			form input[type="text"], form select {
				padding: 10px;
				border: 1px solid #ccc;
				border-radius: 4px;
				width: 100%;
				box-sizing: border-box;
				font-size: 16px;
			}
			.extras {
				display: flex;
				gap: 20px;
				padding: 5px 0;
			}
			.extras label {
				font-size: 1em;
				color: #555;
			}
			.form-row {
				display: flex;
				gap: 15px;
				align-items: flex-end;
			}
			.form-row > * {
				flex-grow: 1;
			}
			.form-row button {
				flex-grow: 0;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<header>
				<h1>☕ Brew & Co - Commandez vos boissons</h1>
			</header>

			<section>
				<h2>Menu</h2>
				<button onclick="fetchMenu()">Charger le menu</button>
				<div id="menu-result"></div>
			</section>

			<section>
				<h2>📝 Passer une commande</h2>
				<form id="order-form">
					<div class="form-row">
						<select name="drink_id" required>
							<option value="">-- Choisir une boisson --</option>
	`)

	// Ajoute les options du menu (utilise la variable 'drinks' de drink.go)
	for _, drink := range drinks {
		fmt.Fprintf(w, `<option value="%s">%s (%.2f€)</option>`, drink.ID, drink.Name, drink.Price)
	}

	// Fin de l'écriture HTML
	fmt.Fprintf(w, `
						</select>
					</div>
					<div class="extras">
						<label><input type="checkbox" name="extras" value="sucre"> Sucre (+0.5€)</label>
						<label><input type="checkbox" name="extras" value="lait"> Lait (+0.5€)</label>
						<label><input type="checkbox" name="extras" value="chantilly"> Chantilly (+0.5€)</label>
					</div>
					<button type="submit">Passer la commande</button>
				</form>
				<div id="order-result"></div>
			</section>

			<section>
				<h2>📋 Voir les commandes</h2>
				<button onclick="fetchOrders()">Charger les commandes</button>
				<div id="orders-result"></div>
			</section>

			<section>
				<h2>✅ Changer le statut d'une commande</h2>
				<form id="update-status-form">
					<div class="form-row">
						<input type="text" name="order_id" placeholder="ID de la commande" required>
						<select name="status" required>
							<option value="">-- Nouveau statut --</option>
							<option value="completed">completed</option>
							<option value="cancelled">cancelled</option>
						</select>
					</div>
					<button type="submit">Mettre à jour</button>
				</form>
				<div id="update-status-result"></div>
			</section>

			<section>
				<h2>❌ Annuler une commande</h2>
				<form id="cancel-order-form">
					<div class="form-row">
						<input type="text" name="order_id" placeholder="ID de la commande" required>
					</div>
					<button type="submit">Annuler</button>
				</form>
				<div id="cancel-order-result"></div>
			</section>
		</div>

		<script>
			// Fonction pour afficher les résultats
			function displayResult(elementId, data) {
				const element = document.getElementById(elementId);
				element.innerHTML = "<pre>" + JSON.stringify(data, null, 2) + "</pre>";
			}

			// Charger le menu
			async function fetchMenu() {
				try {
					const response = await fetch("/menu");
					const data = await response.json();
					displayResult("menu-result", data);
				} catch (error) {
					document.getElementById("menu-result").innerHTML = "<pre>Erreur: " + error.message + "</pre>";
				}
			}

			// Passer une commande
			document.getElementById("order-form").addEventListener("submit", async (e) => {
				e.preventDefault();
				const drinkId = e.target.elements["drink_id"].value;

				// Récupère les extras cochés
				const extras = [];
				document.querySelectorAll('#order-form input[name="extras"]:checked').forEach(checkbox => {
					extras.push(checkbox.value);
				});

				if (!drinkId) {
					alert("Veuillez sélectionner une boisson !");
					return;
				}

				try {
					const response = await fetch("/orders", {
						method: "POST",
						headers: { "Content-Type": "application/json" },
						body: JSON.stringify({
							drink_id: drinkId,
							extras: extras
						})
					});
					const data = await response.json();
					displayResult("order-result", data);
				} catch (error) {
					document.getElementById("order-result").innerHTML = "<pre>Erreur: " + error.message + "</pre>";
				}
			});

			// Charger les commandes
			async function fetchOrders() {
				try {
					const response = await fetch("/orders");
					const data = await response.json();
					displayResult("orders-result", data);
				} catch (error) {
					document.getElementById("orders-result").innerHTML = "<pre>Erreur: " + error.message + "</pre>";
				}
			}

			// Changer le statut d'une commande
			document.getElementById("update-status-form").addEventListener("submit", async (e) => {
				e.preventDefault();
				const orderId = e.target.elements["order_id"].value;
				const status = e.target.elements["status"].value;

				if (!orderId || !status) {
					alert("Veuillez remplir tous les champs !");
					return;
				}

				try {
					const response = await fetch("/orders/" + orderId + "/status", {
						method: "PATCH",
						headers: { "Content-Type": "application/json" },
						body: JSON.stringify({ status: status })
					});
					const data = await response.json();
					displayResult("update-status-result", data);
				} catch (error) {
					document.getElementById("update-status-result").innerHTML = "<pre>Erreur: " + error.message + "</pre>";
				}
			});

			// Annuler une commande
			document.getElementById("cancel-order-form").addEventListener("submit", async (e) => {
				e.preventDefault();
				const orderId = e.target.elements["order_id"].value;

				if (!orderId) {
					alert("Veuillez entrer un ID de commande !");
					return;
				}

				try {
					const response = await fetch("/orders/" + orderId, {
						method: "DELETE"
					});
					const data = await response.json();
					displayResult("cancel-order-result", data);
				} catch (error) {
					document.getElementById("cancel-order-result").innerHTML = "<pre>Erreur: " + error.message + "</pre>";
				}
			});
		</script>
	</body>
	</html>
	`)
}

func main() {
	r := mux.NewRouter()

	// Routes API
	r.HandleFunc("/menu", getMenu).Methods("GET")
	r.HandleFunc("/orders", CreateOrder).Methods("POST")
	r.HandleFunc("/orders", GetOrders).Methods("GET")
	r.HandleFunc("/orders/{id}/status", UpdateOrderStatus).Methods("PATCH")
	r.HandleFunc("/orders/{id}", CancelOrder).Methods("DELETE")

	// Route pour la page de test
	r.HandleFunc("/", testPageHandler).Methods("GET")

	// Démarrage du serveur
	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(r)))
}
