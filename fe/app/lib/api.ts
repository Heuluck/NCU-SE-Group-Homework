import axios from "axios";
import * as mock from "./mockApi";
import type { Task } from "../types";
import { ApiEndpoints } from "../types";

let useMock = typeof window !== "undefined" ? localStorage.getItem("useMockApi") !== "0" : true;

const BASE = (import.meta as any).env?.VITE_API_BASE_URL || "";

const client = axios.create({ baseURL: BASE });

export function isMockEnabled() {
  return useMock;
}

export function setMockEnabled(v: boolean) {
  useMock = v;
  try {
    localStorage.setItem("useMockApi", v ? "1" : "0");
  } catch {}
}

export async function getTasks(): Promise<Task[]> {
  if (useMock) return mock.getTasks();
  const res = await client.get(ApiEndpoints.TASKS);
  return res.data as Task[];
}

export async function addTask(content: string): Promise<Task> {
  if (useMock) return mock.addTask(content);
  const res = await client.post(ApiEndpoints.TASKS, { content });
  return res.data as Task;
}

export async function completeTask(id: number): Promise<Task> {
  if (useMock) return mock.completeTask(id);
  const res = await client.post(ApiEndpoints.COMPLETE(id));
  return res.data as Task;
}

export async function deleteTask(id: number): Promise<void> {
  if (useMock) return mock.deleteTask(id);
  await client.delete(ApiEndpoints.TASK(id));
}

export default {
  isMockEnabled,
  setMockEnabled,
  getTasks,
  addTask,
  completeTask,
  deleteTask,
};
