import {
    createAppContainer,
    createStackNavigator,
    createSwitchNavigator,
} from 'react-navigation';
import React from 'react';
import AsyncStorage from '@react-native-community/async-storage';
import {ActivityIndicator, StatusBar, View} from 'react-native';
import MessageStack from './chat/screens/MessageStack';
import SettingsStack from './chat/screens/SettingsStack';
import {Login} from './auth/screens/Login';
import {Signup} from './auth/screens/Signup';

export class AuthLoadingScreen extends React.Component {
    constructor(props) {
        super(props);
        this._bootstrapAsync();
    }

    // Fetch the token from storage then navigate to our appropriate place
    _bootstrapAsync = async () => {
        const userToken = await AsyncStorage.getItem('userToken');

        // This will switch to the App screen or Auth screen and this loading
        // screen will be unmounted and thrown away.
        this.props.navigation.navigate(userToken ? 'App' : 'Auth');
    };

    // Render any loading content that you like here
    render() {
        return (
            <View>
                <ActivityIndicator />
                <StatusBar barStyle="default" />
            </View>
        );
    }
}

export const ChatNavigator = createStackNavigator({
    Messages: {
        screen: MessageStack,
        navigationOptions: {
            header: null,
        },
    },
    Settings: {
        screen: SettingsStack,
        navigationOptions: {
            header: null,
        },
    },
});

export const AuthNavigator = createStackNavigator(
    {
        Login: {screen: Login},
        Signup: {screen: Signup},
    },
    {
        initialRouteName: 'Login',
    },
);

export default createAppContainer(
    createSwitchNavigator(
        {
            AuthLoading: AuthLoadingScreen,
            App: ChatNavigator,
            Auth: AuthNavigator,
        },
        {
            initialRouteName: 'AuthLoading',
        },
    ),
);
