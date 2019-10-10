import React from 'react';
import {Component} from 'react';
import {View, StatusBar, Text, StyleSheet} from 'react-native';
import AsyncStorage from '@react-native-community/async-storage';
import LinearGradient from 'react-native-linear-gradient';
import {Card, styles} from './Settings';

export default class Edit extends Component {
    static navigationOptions = ({navigation}) => {
        return {
            title: navigation.getParam('name', 'Edit'),
            header: null,
        };
    };

    render() {
        return (
            <LinearGradient colors={['#00af3a', '#005baf']}>
                <View style={styles.container}>
                    <StatusBar barStyle={'light-content'} />
                    <View style={thisStyles.title}>
                        <Text style={styles.title}>{this.props.navigation.getParam('name', 'Edit')}</Text>
                    </View>
                    <Card
                        text={'Back'}
                        textColor={'#444444'}
                        iconName={'arrow-left'}
                        iconType={'feather'}
                        iconColor={'#444444'}
                        showArrow={false}
                        height={50}
                        onPress={() =>
                            this.props.navigation.dispatch({
                                type: 'Navigation/BACK',
                            })
                        }
                    />
                    <View style={styles.spacer} />
                    <Card
                        text={'Name'}
                        textColor={'#444444'}
                        iconName={'user'}
                        iconType={'feather'}
                        iconColor={'#444444'}
                        showArrow={true}
                        height={70}
                        onPress={() =>
                            this.props.navigation.navigate('Edit', {
                                name: 'Name',
                            })
                        }
                    />
                </View>
            </LinearGradient>
        );
    }

    async _signOutAsync() {
        await AsyncStorage.clear();
        this.props.navigation.navigate('Auth');
    }
}

const thisStyles = StyleSheet.create({
    title: {
        paddingTop: 55,
    },
});
