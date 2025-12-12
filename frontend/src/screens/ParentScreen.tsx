import React, { useContext, useState } from 'react';
import { View, Text, StyleSheet, FlatList, TouchableOpacity, TextInput, Alert, SafeAreaView } from 'react-native';
import { GameContext } from '../context/GameContext';

export default function ParentScreen() {
  const { tasks, approveTask, addTask } = useContext(GameContext);
  const [newTaskTitle, setNewTaskTitle] = useState('');
  const [newTaskPoints, setNewTaskPoints] = useState('10');

  const pendingTasks = tasks.filter(t => t.status === 'pending');

  const handleAddTask = () => {
    if (!newTaskTitle.trim()) {
      Alert.alert("提示", "请输入任务名称");
      return;
    }
    addTask(newTaskTitle, parseInt(newTaskPoints));
    setNewTaskTitle('');
    Alert.alert("成功", "任务发布成功！");
  };

  return (
    <SafeAreaView style={styles.container}>
      {/* 待审核区域 */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>🔍 待审核任务 ({pendingTasks.length})</Text>
        <FlatList
          data={pendingTasks}
          keyExtractor={item => item.id.toString()}
          renderItem={({ item }) => (
            <View style={styles.taskItem}>
              <View>
                <Text style={styles.taskTitle}>{item.title}</Text>
                <Text style={styles.taskReward}>奖励: {item.points} 积分</Text>
              </View>
              <TouchableOpacity 
                style={styles.btnApprove} 
                onPress={() => approveTask(item.id)}
              >
                <Text style={styles.btnText}>批准</Text>
              </TouchableOpacity>
            </View>
          )}
          ListEmptyComponent={<Text style={styles.emptyText}>暂无待审核任务</Text>}
        />
      </View>

      {/* 发布任务区域 */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>➕ 发布新任务</Text>
        <View style={styles.formCard}>
          <TextInput 
            style={styles.input}
            placeholder="任务名称 (例如: 背诵英语单词)"
            value={newTaskTitle}
            onChangeText={setNewTaskTitle}
          />
          <View style={styles.row}>
            <Text>积分奖励：</Text>
            <View style={styles.pointsSelector}>
              {['10', '30', '50'].map(p => (
                <TouchableOpacity 
                  key={p}
                  style={[styles.pointOption, newTaskPoints === p && styles.pointOptionActive]}
                  onPress={() => setNewTaskPoints(p)}
                >
                  <Text style={[styles.pointText, newTaskPoints === p && styles.pointTextActive]}>{p}</Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>
          <TouchableOpacity style={styles.btnSubmit} onPress={handleAddTask}>
            <Text style={styles.btnSubmitText}>发布任务</Text>
          </TouchableOpacity>
        </View>
      </View>

      {/* 统计区域 */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>📊 任务概览</Text>
        <View style={styles.statsCard}>
          <Text style={styles.statText}>总任务数: {tasks.length}</Text>
          <Text style={styles.statText}>已完成: {tasks.filter(t => t.status === 'done').length}</Text>
        </View>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#f3f4f6', padding: 16 },
  section: { marginBottom: 24 },
  sectionTitle: { fontSize: 18, fontWeight: 'bold', marginBottom: 12, color: '#1f2937' },
  taskItem: {
    backgroundColor: 'white',
    padding: 16,
    borderRadius: 12,
    marginBottom: 10,
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    borderLeftWidth: 4,
    borderLeftColor: '#d97706'
  },
  taskTitle: { fontSize: 16, fontWeight: '500', color: '#1f2937' },
  taskReward: { color: '#6b7280', fontSize: 12, marginTop: 4 },
  btnApprove: { backgroundColor: '#10b981', paddingVertical: 8, paddingHorizontal: 16, borderRadius: 8 },
  btnText: { color: 'white', fontWeight: 'bold' },
  emptyText: { textAlign: 'center', color: '#9ca3af', padding: 20 },
  
  formCard: { backgroundColor: 'white', padding: 20, borderRadius: 16 },
  input: { borderWidth: 1, borderColor: '#e5e7eb', borderRadius: 8, padding: 12, marginBottom: 16, fontSize: 16 },
  row: { flexDirection: 'row', alignItems: 'center', marginBottom: 20 },
  pointsSelector: { flexDirection: 'row', gap: 10, marginLeft: 10 },
  pointOption: { paddingVertical: 6, paddingHorizontal: 12, borderRadius: 20, backgroundColor: '#f3f4f6' },
  pointOptionActive: { backgroundColor: '#6366f1' },
  pointText: { color: '#6b7280' },
  pointTextActive: { color: 'white', fontWeight: 'bold' },
  btnSubmit: { backgroundColor: '#6366f1', padding: 14, borderRadius: 10, alignItems: 'center' },
  btnSubmitText: { color: 'white', fontWeight: 'bold', fontSize: 16 },

  statsCard: { backgroundColor: 'white', padding: 16, borderRadius: 12 },
  statText: { fontSize: 14, color: '#4b5563', marginBottom: 4 }
});

