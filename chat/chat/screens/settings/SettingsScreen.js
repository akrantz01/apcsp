import React from 'react';
import {Component} from 'react';
import {
    View,
    StatusBar,
    TouchableOpacity,
    Text,
    StyleSheet,
    ScrollView,
} from 'react-native';
import AsyncStorage from '@react-native-community/async-storage';
import LinearGradient from 'react-native-linear-gradient';
import {Icon} from 'react-native-elements';

export class SettingsScreen extends Component {
    static navigationOptions = {
        title: 'Settings',
        header: null,
    };

    render() {
        return (
            <LinearGradient colors={['#00af3a', '#005baf']}>
                <View
                    style={{
                        position: 'relative',
                        width: '100%',
                        height: '100%',
                    }}>
                    <StatusBar barStyle={'dark-content'} />
                    <Text style={styles.title}>Settings</Text>
                    <ScrollView
                        style={{
                            position: 'relative',
                            width: '100%',
                            height: '100%',
                        }}>
                        <Card
                            text={'Back'}
                            textColor={'#444444'}
                            iconName={'arrow-left'}
                            iconType={'feather'}
                            iconColor={'#444444'}
                            showArrow={false}
                            onPress={() =>
                                this.props.navigation.dispatch({
                                    type: 'Navigation/BACK',
                                })
                            }
                        />
                        <View style={{height: 30}} />
                        <View style={{height: 30}} />
                        <Card
                            text={'Sign Out'}
                            textColor={'#ff4444'}
                            iconName={'log-out'}
                            iconType={'feather'}
                            iconColor={'#ff4444'}
                            showArrow={false}
                            onPress={() => this._signOutAsync()}
                        />
                        <View
                            style={{
                                justifyContent: 'center',
                                alignItems: 'center',
                                paddingTop: 10,
                                paddingBottom: 10,
                            }}>
                            <Text style={{color: '#bbbbbb'}}>v0.0.1</Text>
                        </View>
                    </ScrollView>
                </View>
            </LinearGradient>
        );
    }

    async _signOutAsync() {
        await AsyncStorage.clear();
        this.props.navigation.navigate('Auth');
    }
}

class Card extends Component {
    render() {
        return (
            <TouchableOpacity onPress={() => this.props.onPress()}>
                <View style={styles.card}>
                    <LinearGradient
                        start={{x: 0, y: 0}}
                        end={{x: 1, y: 0}}
                        colors={['#ffffff', '#dddddd']}
                        style={styles.grad}>
                        <View
                            style={{
                                position: 'absolute',
                                left: 10,
                            }}>
                            <Icon
                                name={this.props.iconName}
                                type={this.props.iconType}
                                color={this.props.iconColor}
                            />
                        </View>
                        <View
                            style={{
                                position: 'absolute',
                                left: 50,
                            }}>
                            <Text
                                style={[
                                    styles.text,
                                    {color: this.props.textColor},
                                ]}>
                                {this.props.text}
                            </Text>
                        </View>
                        <View
                            style={{
                                position: 'absolute',
                                right: 10,
                            }}>
                            {this.props.showArrow ? (
                                <Icon
                                    name={'chevron-right'}
                                    color={'#444444'}
                                    type={'feather'}
                                />
                            ) : null}
                        </View>
                    </LinearGradient>
                </View>
            </TouchableOpacity>
        );
    }
}

const styles = StyleSheet.create({
    title: {
        marginLeft: 20,
        marginTop: 60,
        marginBottom: 10,
        fontSize: 40,
        color: '#ffffff',
        fontWeight: '800',
    },
    card: {
        marginLeft: 20,
        marginRight: 20,
        marginBottom: 10,
        marginTop: 10,
        borderRadius: 20,
        backgroundColor: '#FFFFFF',
        shadowColor: '#444444',
        shadowOffset: {width: 3, height: 6},
        shadowOpacity: 0.5,
        shadowRadius: 3,
        elevation: 1,
    },
    grad: {
        padding: 10,
        borderRadius: 20,
        flexDirection: 'row',
        flexWrap: 'wrap',
        height: 50,
        justifyContent: 'center',
        alignItems: 'center',
    },
    cardText: {
        padding: 10,
        position: 'relative',
        right: 0,
        top: 0,
        bottom: 0,
        left: 70,
    },
    text: {
        fontSize: 20,
        fontWeight: '600',
    },
});
