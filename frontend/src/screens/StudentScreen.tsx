import React, { useContext } from 'react';
import { View, Text, StyleSheet, FlatList, TouchableOpacity, Alert, SafeAreaView } from 'react-native';
import { GameContext } from '../context/GameContext';

export default function StudentScreen() {
  const { points, tasks, rewards, submitTask, redeemReward } = useContext(GameContext);

  const todoTasks = tasks.filter(t => t.status !== 'done');

  const handleRedeem = (id: number, cost: number, title: string) => {
    Alert.alert("兑换确认", `消耗 ${cost} 积分兑换 "${title}"?`, [
      { text: "取消", style: "cancel" },
      { 
        text: "确定", 
        onPress: () => {
          if (redeemReward(id)) {
            Alert.alert("成功", "兑换成功！快去找家长兑现吧！");
          } else {
            Alert.alert("失败", "积分不足！");
          }
        } 
      }
    ]);
  };

  return (
    <SafeAreaView style={styles.container}>
      {/* 顶部积分卡片 */}
      <View style={styles.headerCard}>
        <Text style={styles.avatar}>🧑‍🚀</Text>
        <View>
          <Text style={styles.welcome}>Hi, 小明同学</Text>
          <Text style={styles.points}>{points} <Text style={{fontSize: 16}}>积分</Text></Text>
        </View>
      </View>

      <View style={styles.section}>
        <Text style={styles.sectionTitle}>📝 今日任务</Text>
        <FlatList
          data={todoTasks}
          keyExtractor={item => item.id.toString()}
          renderItem={({ item }) => (
            <View style={styles.taskItem}>
              <View>
                <Text style={styles.taskTitle}>{item.title}</Text>
                <Text style={styles.taskReward}>+{item.points} 💎</Text>
              </View>
              {item.status === 'todo' ? (
                <TouchableOpacity 
                  style={styles.btnAction} 
                  onPress={() => submitTask(item.id)}
                >
                  <Text style={styles.btnText}>打卡</Text>
                </TouchableOpacity>
              ) : (
                <View style={styles.btnPending}>
                  <Text style={styles.btnTextPending}>审核中</Text>
                </View>
              )}
            </View>
          )}
          ListEmptyComponent={<Text style={styles.emptyText}>任务都完成啦 🎉</Text>}
        />
      </View>

      <View style={[styles.section, { flex: 1 }]}>
        <Text style={styles.sectionTitle}>🎁 奖励兑换</Text>
        <FlatList
          data={rewards}
          numColumns={2}
          keyExtractor={item => item.id.toString()}
          renderItem={({ item }) => (
            <View style={styles.rewardItem}>
              <Text style={{fontSize: 30}}>{item.icon}</Text>
              <Text style={styles.rewardTitle}>{item.title}</Text>
              <Text style={styles.rewardCost}>{item.cost} 积分</Text>
              <TouchableOpacity 
                style={[styles.btnRedeem, points < item.cost && styles.btnDisabled]}
                onPress={() => handleRedeem(item.id, item.cost, item.title)}
                disabled={points < item.cost}
              >
                <Text style={styles.btnText}>兑换</Text>
              </TouchableOpacity>
            </View>
          )}
        />
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#f3f4f6', padding: 16 },
  headerCard: {
    backgroundColor: '#6366f1',
    borderRadius: 20,
    padding: 24,
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 20,
    elevation: 5,
    shadowColor: '#6366f1',
    shadowOffset: {width: 0, height: 4},
    shadowOpacity: 0.3,
    shadowRadius: 5
  },
  avatar: { fontSize: 40, marginRight: 16 },
  welcome: { color: 'rgba(255,255,255,0.8)', fontSize: 14, marginBottom: 4 },
  points: { color: 'white', fontSize: 32, fontWeight: 'bold' },
  section: { marginBottom: 20 },
  sectionTitle: { fontSize: 18, fontWeight: 'bold', marginBottom: 12, color: '#1f2937' },
  taskItem: {
    backgroundColor: 'white',
    padding: 16,
    borderRadius: 12,
    marginBottom: 10,
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center'
  },
  taskTitle: { fontSize: 16, fontWeight: '500', color: '#1f2937' },
  taskReward: { color: '#f59e0b', fontWeight: 'bold', marginTop: 4 },
  btnAction: { backgroundColor: '#6366f1', paddingVertical: 8, paddingHorizontal: 16, borderRadius: 20 },
  btnPending: { backgroundColor: '#fef3c7', paddingVertical: 8, paddingHorizontal: 16, borderRadius: 20 },
  btnText: { color: 'white', fontWeight: 'bold', fontSize: 12 },
  btnTextPending: { color: '#d97706', fontWeight: 'bold', fontSize: 12 },
  emptyText: { textAlign: 'center', color: '#9ca3af', marginTop: 20 },
  rewardItem: {
    flex: 1,
    backgroundColor: 'white',
    margin: 5,
    borderRadius: 12,
    padding: 16,
    alignItems: 'center',
    gap: 8
  },
  rewardTitle: { fontSize: 14, fontWeight: '500', textAlign: 'center' },
  rewardCost: { color: '#f59e0b', fontWeight: 'bold', fontSize: 12 },
  btnRedeem: { backgroundColor: '#ec4899', paddingVertical: 6, paddingHorizontal: 20, borderRadius: 20, marginTop: 4 },
  btnDisabled: { backgroundColor: '#e5e7eb' }
});

