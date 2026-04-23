import axios from "axios";
import type { Task } from "../types";
import { ApiEndpoints } from "../types";

const PRISM_DEFAULT = "http://127.0.0.1:8080";
const env = (import.meta as any).env || {};
const isDev = Boolean(env.DEV);
const envBase = env.VITE_API_BASE_URL || "";
const BASE = envBase || (isDev ? PRISM_DEFAULT : "");

const client = axios.create({ baseURL: BASE });

// 统一响应拦截器：将后端返回的错误信息封装成 Error，便于上层组件统一展示
client.interceptors.response.use(
  (response) => response,
  (error) => {
    const res = error?.response;
    const msg = res?.data?.message || error?.message || "网络错误";
    const apiError: any = new Error(msg);
    apiError.status = res?.status;
    apiError.data = res?.data;
    return Promise.reject(apiError);
  }
);

export async function getTasks(): Promise<Task[]> {
  const res = await client.get(ApiEndpoints.TASKS);
  return res.data as Task[];
}

export type FetchTasksOptions = {
  status?: string; // 'pending' | 'completed' | 'all'
  page?: number;
  size?: number;
};

// Client-side filtering + pagination helper. Backend mock returns full array,
// so we filter/slice on the frontend.
export async function fetchTasks(opts?: FetchTasksOptions): Promise<{ items: Task[]; total: number }> {
  const res = await client.get(ApiEndpoints.TASKS);
  let data = res.data as Task[];

  if (opts?.status && opts.status !== "all") {
    data = data.filter((t) => t.status === opts.status);
  }

  const total = data.length;

  if (opts?.page && opts?.size) {
    const start = (opts.page - 1) * opts.size;
    data = data.slice(start, start + opts.size);
  }

  return { items: data, total };
}

export async function addTask(content: string): Promise<Task> {
  const res = await client.post(ApiEndpoints.TASKS, { content });
  return res.data as Task;
}

export async function completeTask(id: number): Promise<Task> {
  const res = await client.post(ApiEndpoints.COMPLETE(id));
  return res.data as Task;
}

export async function deleteTask(id: number): Promise<void> {
  await client.delete(ApiEndpoints.TASK(id));
}

export default {
  getTasks,
  fetchTasks,
  addTask,
  completeTask,
  deleteTask,
};
