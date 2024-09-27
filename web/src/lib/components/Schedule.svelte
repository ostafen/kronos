<script>
  export let info;
  import ScheduleService from "../ScheduleService";
  import { schedules } from "../store";
  import Swal from "sweetalert2";
  const scheduleService = new ScheduleService("http://localhost:9175");
  import cronstrue from "cronstrue";

  let dialog;
  function openDialog() {
    dialog.showModal();
  }

  function closeDialog() {
    dialog.close();
  }

  const pauseTask = async () => {
    let res = await scheduleService.pauseSchedule(info.id);

    let scheduleList = $schedules;
    let index = scheduleList.findIndex((item) => item.id === info.id);
    scheduleList[index].status = "paused";
    schedules.set(scheduleList);

    Swal.fire("Task paused", "", "info");
  };

  const resumeTask = async () => {
    let res = await scheduleService.resumeSchedule(info.id);

    let scheduleList = $schedules;
    let index = scheduleList.findIndex((item) => item.id === info.id);
    scheduleList[index].status = "active";
    schedules.set(scheduleList);

    Swal.fire("Task resumed", "", "info");
  };

  const triggerTask = async () => {
    let res = await scheduleService.triggerSchedule(info.id);

    let scheduleList = $schedules;
    let index = scheduleList.findIndex((item) => item.id === info.id);
    scheduleList[index].status = "running";
    schedules.set(scheduleList);

    Swal.fire("Task triggered", "", "info");
  };

  const deleteTask = async () => {
    Swal.fire({
      title: "Do you want to delete this task?",
      showDenyButton: true,
      showConfirmButton: false,
      denyButtonText: `Delete`,
    }).then(async (result) => {
      /* Read more about isConfirmed, isDenied below */
      if (result.isConfirmed) {
        Swal.fire(JSON.stringify(result), "", "success");
      } else if (result.isDenied) {
        let res = await scheduleService.deleteSchedule(info.id);

        let scheduleList = $schedules;

        let index = scheduleList.findIndex((item) => item.id === info.id);
        scheduleList.splice(index, 1);
        schedules.set(scheduleList);

        //let index = scheduleList.findIndex((item) => item.id === info.id);
        // scheduleList.splice(index, 1);
        // schedules.set(scheduleList);

        Swal.fire("Task deleted", "", "info");
      }
    });
  };
</script>

<tr>
  <td>{info.title}</td>
  <td>{info.description}</td>
  <td><a class="truncate" href={info.url}>{info.url}</a></td>
  <td>
    {#if info.isRecurring}
      <!-- content here -->
      <div class="badge badge-primary">Recurring</div>
    {:else}
      <div class="badge badge-secondary">Run Once</div>

      <!-- else content here -->
    {/if}
  </td>
  <td
    ><div
      class="badge badge-{info.status === 'active' ? 'success' : 'warning'}"
    >
      {info.status}
    </div></td
  >
  <td>
    <div
      class="tooltip"
      data-tip={info.status === "active" ? "Pause" : "Resume"}
    >
      <button
        on:click={() => (info.status === "active" ? pauseTask() : resumeTask())}
        class="btn btn-sm shadow-sm btn-{info.status === 'active'
          ? 'secondary'
          : 'success'}"
      >
        {#if info.status === "active"}
          <!-- content here --><i class="bi bi-pause"></i>
        {:else}
          <i class="bi bi-play"></i>

          <!-- else content here -->
        {/if}
      </button>
    </div>
    <div class="tooltip" data-tip="Delete Task">
      <button class="btn btn-sm btn-error" on:click={() => deleteTask()}
        ><i class="bi bi-trash"></i></button
      >
    </div>
    <div class="tooltip" data-tip="More Info">
      <button class="btn btn-sm bg-blue-300" on:click={openDialog}
        ><i class="bi bi-eye"></i></button
      >
    </div>
    <div class="tooltip" data-tip="Trigger">
      <button class="btn btn-sm bg-black text-white" on:click={triggerTask}
        ><i class="bi bi-lightning"></i></button
      >
    </div>
  </td>
</tr>
<dialog bind:this={dialog} class="modal">
  <div class="modal-box">
    <div class="task">
      <h3 class="text-lg font-bold">{info.title}!</h3>

      <div><span class="key">Status:</span> {info.status}</div>
      <div><span class="key my-4">Description:</span> {info.description}</div>
      <div><span class="key">ID:</span> {info.id}</div>
      <div>
        <span class="key">URL:</span>
        <a href={info.url} target="_blank">{info.url}</a>
      </div>
      <hr />
      <div>
        <span class="key">Created At:</span>
        {new Date(info.createdAt).toLocaleString()}
      </div>
      {#if info.isRecurring}
        <div><span class="key">Cron Expression:</span> {info.cronExpr}</div>
        <div>
          <span class="key">Cron Expression:</span>
          {cronstrue.toString(info.cronExpr)}
        </div>
      {:else}
        <div>
          <span class="key">Run At:</span>
          {new Date(info.runAt).toLocaleString()}
        </div>
      {/if}
      <div>
        <span class="key">Start At:</span>
        {new Date(info.startAt).toLocaleString()}
      </div>
      <div>
        <span class="key">End At:</span>
        {new Date(info.endAt).toLocaleString()}
      </div>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>

<style>
  .task {
    margin-bottom: 1em;
  }
  .key {
    font-weight: bold;
    margin-top: 30px;
  }
</style>
