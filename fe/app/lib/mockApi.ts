import { type Task } from "../types";

let nextId = 3;
let tasks: Task[] = [
  { id: 1, content: "示例任务：阅读规格文档", status: "pending", created_at: new Date().toISOString() },
  { id: 2, content: "示例任务：实现前端 Mock API", status: "completed", created_at: new Date().toISOString(), completed_at: new Date().toISOString() },
];

const delay = (ms = 200) => new Promise((res) => setTimeout(res, ms));

export async function getTasks(): Promise<Task[]> {
  await delay();
  return tasks.slice().sort((a, b) => a.id - b.id);
}

export async function addTask(content: string): Promise<Task> {
  await delay();
  const t: Task = { id: nextId++, content, status: "pending", created_at: new Date().toISOString() };
  tasks.push(t);
  return t;
}

export async function completeTask(id: number): Promise<Task> {
  await delay();
  const t = tasks.find((x) => x.id === id);
  if (!t) throw new Error("not_found");
  if (t.status === "completed") return t;
  t.status = "completed";
  t.completed_at = new Date().toISOString();
  return t;
}

export async function deleteTask(id: number): Promise<void> {
  await delay();
  const idx = tasks.findIndex((x) => x.id === id);
  if (idx === -1) throw new Error("not_found");
  tasks.splice(idx, 1);
}

export async function resetMockData() {
  tasks = [
    { id: 1, content: "示例任务：阅读规格文档", status: "pending", created_at: new Date().toISOString() },
    { id: 2, content: "示例任务：实现前端 Mock API", status: "completed", created_at: new Date().toISOString(), completed_at: new Date().toISOString() },
  ];
  nextId = 3;
}
