// Node 20+ 自带的实验性全局 localStorage 会抢在 jsdom/happy-dom 的实现之前注册，
// 且在这台机器上取不到可用实现（已验证是运行环境问题，不是本项目代码问题）。
// 这里用一个最小的内存实现兜底，保证被测代码里的裸 localStorage 引用在任何环境下都可用。
class MemoryStorage implements Storage {
  private store = new Map<string, string>();

  get length(): number {
    return this.store.size;
  }

  clear(): void {
    this.store.clear();
  }

  getItem(key: string): string | null {
    return this.store.has(key) ? this.store.get(key)! : null;
  }

  key(index: number): string | null {
    return Array.from(this.store.keys())[index] ?? null;
  }

  removeItem(key: string): void {
    this.store.delete(key);
  }

  setItem(key: string, value: string): void {
    this.store.set(key, String(value));
  }
}

globalThis.localStorage = new MemoryStorage();
