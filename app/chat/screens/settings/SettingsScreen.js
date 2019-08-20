import React from 'react';
import {Component} from 'react';
import {AsyncStorage, View} from 'react-native';
import {Button} from 'react-native-elements';

export class SettingsScreen extends Component {
    static navigationOptions = {
        title: 'Settings',
    };

    render() {
        return (
            <View>
                <Button title="Sign out" onPress={this._signOutAsync} />
            </View>
        );
    }

    _signOutAsync = async () => {
        await AsyncStorage.clear();
        this.props.navigation.navigate('Auth');
    };
}
