import deleteAllSchedules from './delete.js';
import schedules from './schedules.json';

export default async function seedDatabase() {
  try {
    await deleteAllSchedules();
    const settledPromises = await Promise.allSettled(
      schedules.map(async (schedule) => {
        const response = await fetch('http://localhost:9175/api/v1/schedules', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(schedule),
        });

        return response.json();
      })
    );

    console.log('Database seeded');
    console.log(settledPromises);
  } catch (err) {
    console.error('There was an error during the seeding', err);
  }
}

if (typeof process !== 'undefined' && process.argv[2] === 'exec') {
  void seedDatabase();
}
