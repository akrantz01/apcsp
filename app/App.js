import React from 'react';
import {createBottomTabNavigator, createAppContainer} from 'react-navigation';
import Icon from 'react-native-vector-icons/Feather';
import HomeScreen from './screens/HomeScreen';
import SettingsScreen from './screens/SettingsScreen';

const MainNavigator = createBottomTabNavigator(
    {
        Messages: {screen: HomeScreen},
        Settings: {screen: SettingsScreen},
    },
    {
        defaultNavigationOptions: ({navigation}) => ({
            tabBarIcon: ({tintColor}) => {
                const {routeName} = navigation.state;
                let iconName;
                if (routeName === 'Messages') {
                    iconName = 'message-circle';
                } else if (routeName === 'Settings') {
                    iconName = 'settings';
                }

                // You can return any component that you like here!
                return <Icon name={iconName} size={25} color={tintColor} />;
            },
        }),
    },
);

const App = createAppContainer(MainNavigator);

export default App;
