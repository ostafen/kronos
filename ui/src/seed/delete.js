export default async function deleteAllSchedules() {
    try {
        const schedules = await fetch("http://localhost:9175/api/v1/schedules").then(
            (d) => d.json(),
        );

        const settledPromises = await Promise.allSettled(
            schedules.map(schedule =>
                fetch(`http://localhost:9175/api/v1/schedules/${schedule.id}`, {
                    method: "DELETE",
                }))
        );

        console.log("Schedules deleted");
        console.log(JSON.stringify(settledPromises, null, 4));
    } catch (error) {
        console.error("There was an error!");
        console.error(error);
    }
}

if (typeof process !== "undefined" && process.argv[2] === 'exec') {
    void deleteAllSchedules();
}