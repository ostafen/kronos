export default async function logAllSchedules() {
  const response = await fetch('http://localhost:9175/api/v1/schedules');
  console.log(await response.json());
}

if (typeof process !== 'undefined' && process.argv[2] === 'exec') {
  void logAllSchedules();
}
