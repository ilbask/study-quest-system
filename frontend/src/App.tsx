import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { GameProvider } from './context/GameContext';
import StudentScreen from './screens/StudentScreen';
import ParentScreen from './screens/ParentScreen';
import { Text } from 'react-native';

const Tab = createBottomTabNavigator();

export default function App() {
  return (
    <GameProvider>
      <NavigationContainer>
        <Tab.Navigator
          screenOptions={({ route }) => ({
            headerShown: false,
            tabBarIcon: ({ focused, color, size }) => {
              let icon = '❓';
              if (route.name === '学生端') icon = '👦';
              if (route.name === '家长端') icon = '🛡️';
              return <Text style={{fontSize: size}}>{icon}</Text>;
            },
            tabBarActiveTintColor: '#6366f1',
            tabBarInactiveTintColor: 'gray',
          })}
        >
          <Tab.Screen name="学生端" component={StudentScreen} />
          <Tab.Screen name="家长端" component={ParentScreen} />
        </Tab.Navigator>
      </NavigationContainer>
    </GameProvider>
  );
}

