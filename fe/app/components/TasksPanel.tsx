import React, { useEffect, useState } from "react";
import { Button, Input, Space, Tag, Popconfirm, message, Typography, Card, Table } from "antd";
import type { ColumnsType } from "antd/es/table";
import api from "../lib/api";
import type { Task } from "../types";
import TaskFilter from "./TaskFilter";
import Pagination from "./Pagination";
import { showError } from "./ErrorToast";

const { Text } = Typography;

export function TasksPanel() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(false);
  const [value, setValue] = useState("");
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(5);
  const [filter, setFilter] = useState<string>("all");
  const [total, setTotal] = useState(0);

  useEffect(() => {
    refresh();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filter, page, pageSize]);

  async function refresh() {
    setLoading(true);
    try {
      const res = await api.fetchTasks({ status: filter, page, size: pageSize });
      setTasks(res.items);
      setTotal(res.total);
    } catch (err) {
      console.error(err);
      showError(err);
    } finally {
      setLoading(false);
    }
  }

  async function handleAdd() {
    if (!value.trim()) {
      message.warning("请输入任务内容");
      return;
    }
    try {
      const t = await api.addTask(value.trim());
      setValue("");
      message.success("已添加任务");
      // refresh to reflect filter/pagination
      setPage(1);
      await refresh();
    } catch (err) {
      console.error(err);
      showError(err);
    }
  }

  async function handleComplete(id: number) {
    try {
      const t = await api.completeTask(id);
      // refresh current page
      await refresh();
      message.success("标记已完成");
    } catch (err) {
      console.error(err);
      showError(err);
    }
  }

  async function handleDelete(id: number) {
    try {
      await api.deleteTask(id);
      // refresh current page
      await refresh();
      message.success("已删除");
    } catch (err) {
      console.error(err);
      showError(err);
    }
  }

  // mock switch removed — always use configured API (Prism in dev)

  return (
    <Card style={{ width: "100%", maxWidth: 900 }}>
      <Space direction="vertical" style={{ width: "100%" }}>
        <Space style={{ justifyContent: "space-between", width: "100%" }}>
          <Space>
            <Text strong>任务面板</Text>
            <Text type="secondary">(与后端 API 对接)</Text>
          </Space>
        </Space>

        <Space style={{ width: "100%", justifyContent: "space-between" }}>
          <Space style={{ flex: 1 }}>
            <Input
              placeholder="新任务内容，按回车或点击添加"
              value={value}
              onChange={(e) => setValue(e.target.value)}
              onPressEnter={handleAdd}
            />
            <Button type="primary" onClick={handleAdd}>
              添加
            </Button>
            <Button onClick={refresh}>刷新</Button>
          </Space>
          <Space>
            <TaskFilter value={filter} onChange={(v) => { setFilter(v); setPage(1); }} />
          </Space>
        </Space>

        <Table<Task>
          rowKey="id"
          loading={loading}
          dataSource={tasks}
          pagination={false}
          columns={(
            (() => {
              const cols: ColumnsType<Task> = [
                {
                  title: "内容",
                  dataIndex: "content",
                  key: "content",
                  render: (_text, record) => (
                    <div>
                      <div>{record.content}</div>
                      <div style={{ color: "rgba(0,0,0,0.45)", marginTop: 6 }}>
                        <div>创建于：{new Date(record.created_at).toLocaleString()}</div>
                        {record.completed_at && <div>完成于：{new Date(record.completed_at).toLocaleString()}</div>}
                      </div>
                    </div>
                  ),
                },
                {
                  title: "状态",
                  dataIndex: "status",
                  key: "status",
                  width: 120,
                  render: (status: string) => (status === "completed" ? <Tag color="green">已完成</Tag> : <Tag>待办</Tag>),
                },
                {
                  title: "操作",
                  key: "actions",
                  width: 160,
                  render: (_text, record) => (
                    <Space>
                      {record.status !== "completed" ? (
                        <Button type="link" onClick={() => handleComplete(record.id)}>
                          完成
                        </Button>
                      ) : null}
                      <Popconfirm title="确认删除？" onConfirm={() => handleDelete(record.id)}>
                        <Button type="link" danger>
                          删除
                        </Button>
                      </Popconfirm>
                    </Space>
                  ),
                },
              ];
              return cols;
            })()
          )}
        />
        <Pagination
          current={page}
          pageSize={pageSize}
          total={total}
          onChange={async (p, size) => {
            setPage(p);
            if (size && size !== pageSize) setPageSize(size);
            await refresh();
          }}
        />
      </Space>
    </Card>
  );
}

export default TasksPanel;
