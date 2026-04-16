export const TaskFields = {
  ID: "id",
  CONTENT: "content",
  STATUS: "status",
  CREATED_AT: "created_at",
  COMPLETED_AT: "completed_at",
} as const;

export interface Task {
  id: number;
  content: string;
  status: string; // 'pending' | 'completed'
  created_at: string;
  completed_at?: string | null;
}

export type TaskInput = {
  content: string;
};

export const ApiEndpoints = {
  TASKS: "/tasks",
  TASK: (id: number | string) => `/tasks/${id}`,
  COMPLETE: (id: number | string) => `/tasks/${id}/complete`,
};
