<script>
  import Schedule from "./Schedule.svelte";
  import ScheduleService from "../ScheduleService";
  import { schedules } from "../store";
  import { onMount } from "svelte";

  let scheduleList = [];
  schedules.subscribe((value) => {
    scheduleList = value;
  });

  const scheduleService = new ScheduleService("http://localhost:9175");
  const loadData = async () => {
    let res = await scheduleService.getAllSchedules();
    schedules.set(res);
  };

  loadData();
</script>

<div class="w-4/5">
  <table class="table w-full table-zebra">
    <thead>
      <tr>
        <th>Name</th>
        <th>Description</th>
        <th>URL</th>
        <th>Type</th>
        <th>Status</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {#each $schedules as item}
        <Schedule info={item} />
      {/each}
    </tbody>
  </table>
</div>
