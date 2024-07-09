<script>
  let isRecurring = true;
  import ScheduleService from "$lib/ScheduleService";
  import { schedules } from "$lib/store";
  const scheduleService = new ScheduleService("http://localhost:9175");
  let formData = {};
  async function handleSubmit(event) {
    event.preventDefault();
    const form = event.target;
    const data = new FormData(form);

    form.reset();
    formData = Object.fromEntries(data.entries());
    if (!formData.isRecurring) {
      formData.isRecurring = false;
    } else {
      formData.isRecurring = true;
    }
    if (formData.startAt) {
      formData.startAt = new Date(formData.startAt).toISOString();
    }
    if (formData.endAt) {
      formData.endAt = new Date(formData.endAt).toISOString();
    }

    if (formData.runTime && formData.runDate) {
      formData.runAt = new Date(
        `${formData.runDate}T${formData.runTime}:00`
      ).toISOString();

      formData.startAt = formData.runAt;
      formData.endAt = formData.runAt;
      delete formData.runTime;
      delete formData.runDate;
    }

    let result = await scheduleService.registerSchedule(formData);
    let scheduleList = $schedules;
    scheduleList.push(result);
    $schedules = scheduleList;

    // Add your form submission logic here (e.g., send data to a server)
  }
</script>

<button class="btn" onclick="my_modal_2.showModal()">Create New Task</button>
<dialog id="my_modal_2" class="modal">
  <div class="modal-box">
    <div
      class=" bg-card text-card-foreground w-full max-w-2xl"
      data-v0-t="card"
    >
      <div class="flex flex-col space-y-1.5 p-6">
        <h3
          class="whitespace-nowrap text-2xl font-semibold leading-none tracking-tight"
        >
          Create New Task
        </h3>
        <p class="text-sm text-muted-foreground">
          Set up a new schedule for your application.
        </p>
      </div>
      <form method="POST" on:submit={handleSubmit}>
        <div class="p-6 grid gap-6">
          <div class="grid gap-4">
            <div class="grid gap-2">
              <label
                class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                for="title"
              >
                Title
              </label>
              <input
                class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                id="title"
                name="title"
                placeholder="Enter a title"
                required
              />
            </div>
            <div class="grid gap-2">
              <label
                class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                for="description"
              >
                Description
              </label>
              <textarea
                class="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                id="description"
                name="description"
                placeholder="Describe the schedule"
                required
              ></textarea>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div class="grid gap-2 items-center">
                <div class="form-control">
                  <span class="label-text">Is Recurring</span>

                  <input
                    id="isRecurring"
                    name="isRecurring"
                    class="toggle"
                    type="checkbox"
                    bind:checked={isRecurring}
                  />
                </div>
              </div>
              <div class="grid gap-2">
                {#if isRecurring}
                  <!-- content here -->
                  <label
                    class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                    for="cron-expression"
                  >
                    Cron Expression
                  </label>
                  <input
                    class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    id="cronExpr"
                    name="cronExpr"
                    required
                    placeholder="Enter a cron expression"
                  />
                {:else}
                  <!-- else content here -->
                  <label
                    class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                    for="cron-expression"
                  >
                    Run At
                  </label>
                  <input
                    type="time"
                    required
                    id="runTime"
                    name="runTime"
                    class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  />
                  <input
                    type="date"
                    required
                    id="runTime"
                    name="runDate"
                    class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  />
                {/if}
              </div>
            </div>
            <div class="grid gap-2">
              <label
                class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                for="webhook-url"
              >
                Webhook URL
              </label>
              <input
                class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                id="url"
                name="url"
                placeholder="Enter a webhook URL"
                required=""
                type="url"
              />
            </div>
            {#if isRecurring}
              <div class="grid grid-cols-2 gap-4">
                <div class="grid gap-2">
                  <label
                    class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                    for="run-at"
                  >
                    Start On
                  </label><input
                    type="date"
                    id="startAt"
                    name="startAt"
                    class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  />
                </div>
                <div class="grid gap-2">
                  <label
                    class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                    for="start-at"
                  >
                    End On
                  </label>
                  <input
                    type="date"
                    id="endAt"
                    required
                    name="endAt"
                    class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  />
                </div>
              </div>
              <!-- content here -->
            {/if}
            <div class="grid gap-4">
              <div class="grid gap-2">
                <label
                  class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                  for="metadata"
                >
                  Metadata
                </label>
                <textarea
                  class="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 min-h-[100px]"
                  id="metadata"
                  placeholder="Enter a JSON object"
                ></textarea>
              </div>
            </div>
          </div>
        </div>
        <div class="flex items-center p-6">
          <button
            class="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2 ml-auto"
            type="submit"
          >
            Create Schedule
          </button>
        </div>
      </form>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop">
    <button>close</button>
  </form>
</dialog>
