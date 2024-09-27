class ScheduleService {
  constructor(baseURL) {
    this.baseURL = baseURL;
  }

  async registerSchedule(data) {
    try {
      const response = await fetch(`${this.baseURL}/schedules`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error(`Error: ${response.status} ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      throw error;
    }
  }

  async getSchedule(id) {
    try {
      const response = await fetch(`${this.baseURL}/schedules/${id}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`Error: ${response.status} ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      throw error;
    }
  }

  async getAllSchedules() {
    try {
      const response = await fetch(`${this.baseURL}/schedules`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`Error: ${response.status} ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      throw error;
    }
  }

  async deleteSchedule(id) {
    try {
      const response = await fetch(`${this.baseURL}/schedules/${id}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`Error: ${response.status} ${response.statusText}`);
      }

      return await { status: response.status, message: "Schedule deleted" }; //await response.text();
    } catch (error) {
      throw error;
    }
  }

  async pauseSchedule(id) {
    try {
      const response = await fetch(`${this.baseURL}/schedules/${id}/pause`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`Error: ${response.status} ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      throw error;
    }
  }

  async resumeSchedule(id) {
    try {
      const response = await fetch(`${this.baseURL}/schedules/${id}/resume`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`Error: ${response.status} ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      throw error;
    }
  }

  async triggerSchedule(id) {
    try {
      const response = await fetch(`${this.baseURL}/schedules/${id}/trigger`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`Error: ${response.status} ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      throw error;
    }
  }
}

export default ScheduleService;
