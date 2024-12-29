import fs from "fs";

// Leggi il file JSON contenente gli schedule
fs.readFile("schedules.json", "utf8", (err, data) => {
  if (err) {
    console.error("Errore nel leggere il file JSON:", err);
    return;
  }

  // Parsea il JSON
  const schedules = JSON.parse(data);

  // Itera su ogni oggetto JSON
  schedules.forEach((schedule) => {
    const title = schedule.title;
    const description = schedule.description;
    const url = schedule.url;
    const isRecurring = schedule.isRecurring;
    const startAt = schedule.startAt;
    const endAt = schedule.endAt;
    const runAt = schedule.runAt;
    const cronExpr = schedule.cronExpr;

    // Crea il corpo della richiesta JSON
    const requestData = {
      title: title,
      description: description,
      url: url,
      isRecurring: isRecurring,
      startAt: startAt,
      endAt: endAt,
      runAt: runAt,
      cronExpr: cronExpr,
    };

    // Esegui la richiesta POST con fetch
    fetch("http://localhost:9175/api/v1/schedules", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(requestData),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log(`Risposta per "${title}":`, data);
      })
      .catch((error) => {
        console.error("Errore nella richiesta:", error);
      });
  });
});
