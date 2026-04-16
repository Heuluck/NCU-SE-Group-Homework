import axios from "axios";
import type { Task } from "../types";
import { ApiEndpoints } from "../types";

const PRISM_DEFAULT = "http://127.0.0.1:4010";
const env = (import.meta as any).env || {};
const isDev = Boolean(env.DEV);
const envBase = env.VITE_API_BASE_URL || "";
const BASE = envBase || (isDev ? PRISM_DEFAULT : "");

const client = axios.create({ baseURL: BASE });

export async function getTasks(): Promise<Task[]> {
  const res = await client.get(ApiEndpoints.TASKS);
  return res.data as Task[];
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
  addTask,
  completeTask,
  deleteTask,
};
