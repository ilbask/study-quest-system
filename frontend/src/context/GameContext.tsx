import React, { createContext, useState, useEffect, ReactNode } from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';

// 类型定义
export type TaskStatus = 'todo' | 'pending' | 'done';

export interface Task {
  id: number;
  title: string;
  points: number;
  status: TaskStatus;
}

export interface Reward {
  id: number;
  title: string;
  cost: number;
  icon: string;
}

interface GameContextType {
  points: number;
  tasks: Task[];
  rewards: Reward[];
  submitTask: (id: number) => void;
  approveTask: (id: number) => void;
  addTask: (title: string, points: number) => void;
  redeemReward: (id: number) => boolean;
}

export const GameContext = createContext<GameContextType>({} as GameContextType);

export const GameProvider = ({ children }: { children: ReactNode }) => {
  const [points, setPoints] = useState(100);
  const [tasks, setTasks] = useState<Task[]>([
    { id: 1, title: '完成数学口算', points: 30, status: 'todo' },
    { id: 2, title: '整理书包', points: 10, status: 'todo' },
  ]);
  const [rewards, setRewards] = useState<Reward[]>([
    { id: 1, title: '看电视 30分钟', cost: 50, icon: '📺' },
    { id: 2, title: '玩手机 15分钟', cost: 30, icon: '📱' },
  ]);

  // 加载数据
  useEffect(() => {
    const loadData = async () => {
      try {
        const savedPoints = await AsyncStorage.getItem('points');
        const savedTasks = await AsyncStorage.getItem('tasks');
        if (savedPoints) setPoints(parseInt(savedPoints));
        if (savedTasks) setTasks(JSON.parse(savedTasks));
      } catch (e) {
        console.error("Failed to load data", e);
      }
    };
    loadData();
  }, []);

  // 保存数据
  const saveData = async (newPoints: number, newTasks: Task[]) => {
    try {
      await AsyncStorage.setItem('points', newPoints.toString());
      await AsyncStorage.setItem('tasks', JSON.stringify(newTasks));
    } catch (e) {
      console.error("Failed to save data", e);
    }
  };

  const submitTask = (id: number) => {
    const newTasks = tasks.map(t => t.id === id ? { ...t, status: 'pending' as TaskStatus } : t);
    setTasks(newTasks);
    saveData(points, newTasks);
  };

  const approveTask = (id: number) => {
    const task = tasks.find(t => t.id === id);
    if (task && task.status !== 'done') {
      const newPoints = points + task.points;
      const newTasks = tasks.map(t => t.id === id ? { ...t, status: 'done' as TaskStatus } : t);
      setPoints(newPoints);
      setTasks(newTasks);
      saveData(newPoints, newTasks);
    }
  };

  const addTask = (title: string, pointsVal: number) => {
    const newTask: Task = { id: Date.now(), title, points: pointsVal, status: 'todo' };
    const newTasks = [newTask, ...tasks];
    setTasks(newTasks);
    saveData(points, newTasks);
  };

  const redeemReward = (id: number): boolean => {
    const reward = rewards.find(r => r.id === id);
    if (reward && points >= reward.cost) {
      const newPoints = points - reward.cost;
      setPoints(newPoints);
      saveData(newPoints, tasks);
      return true;
    }
    return false;
  };

  return (
    <GameContext.Provider value={{ points, tasks, rewards, submitTask, approveTask, addTask, redeemReward }}>
      {children}
    </GameContext.Provider>
  );
};

