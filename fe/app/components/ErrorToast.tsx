import { message } from "antd";

export function showError(err: any) {
  const msg = err?.message || err?.data?.message || "操作失败";
  message.error(msg);
}

export default showError;
