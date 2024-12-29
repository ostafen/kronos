try {
  const schedules = await fetch("http://localhost:9175/api/v1/schedules").then(
    (d) => d.json(),
  );

  for (const schedule of schedules) {
    await fetch(`http://localhost:9175/api/v1/schedules/${schedule.id}`, {
      method: "DELETE",
    });
  }

  console.log("All schedules were deleted successfully!");
} catch (error) {
  console.error("There was an error!");
  console.error(error);
}
