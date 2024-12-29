const response = await fetch("http://localhost:9175/api/v1/schedules");

const json = await response.json();

console.log(JSON.stringify(json, null, 4));
