import React, { useEffect, useState } from "react";
import { Button, Input, Space, Switch, Tag, Popconfirm, message, Typography, Card, Table } from "antd";
import type { ColumnsType } from "antd/es/table";
import api from "../lib/api";
import type { Task } from "../types";

const { Text } = Typography;

export function TasksPanel() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(false);
  const [value, setValue] = useState("");
  const [mockEnabled, setMockEnabled] = useState<boolean>(api.isMockEnabled());

  useEffect(() => {
    refresh();
  }, []);

  async function refresh() {
    setLoading(true);
    try {
      const res = await api.getTasks();
      setTasks(res);
    } catch (err) {
      console.error(err);
      message.error("加载任务失败");
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
      setTasks((s) => [...s, t]);
    } catch (err) {
      console.error(err);
      message.error("添加失败");
    }
  }

  async function handleComplete(id: number) {
    try {
      const t = await api.completeTask(id);
      setTasks((s) => s.map((x) => (x.id === id ? t : x)));
      message.success("标记已完成");
    } catch (err) {
      console.error(err);
      message.error("标记失败");
    }
  }

  async function handleDelete(id: number) {
    try {
      await api.deleteTask(id);
      setTasks((s) => s.filter((x) => x.id !== id));
      message.success("已删除");
    } catch (err) {
      console.error(err);
      message.error("删除失败");
    }
  }

  function toggleMock(v: boolean) {
    api.setMockEnabled(v);
    setMockEnabled(v);
    message.info(v ? "已启用 Mock 模式" : "已切换为真实 API 模式");
    // reload tasks from selected backend
    refresh();
  }

  return (
    <Card style={{ width: "100%", maxWidth: 900 }}>
      <Space direction="vertical" style={{ width: "100%" }}>
        <Space style={{ justifyContent: "space-between", width: "100%" }}>
          <Space>
            <Text strong>任务面板</Text>
            <Text type="secondary">(支持 Mock / 实时 API 切换)</Text>
          </Space>
          <Space>
            <Text>Mock</Text>
            <Switch checked={mockEnabled} onChange={toggleMock} />
          </Space>
        </Space>

        <Space style={{ width: "100%" }}>
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
      </Space>
    </Card>
  );
}

export default TasksPanel;
