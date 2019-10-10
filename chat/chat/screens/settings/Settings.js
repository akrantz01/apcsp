import React from 'react';
import {Component} from 'react';
import {View, StatusBar, TouchableOpacity, Text, StyleSheet, ScrollView} from 'react-native';
import LinearGradient from 'react-native-linear-gradient';
import {Icon} from 'react-native-elements';
import {APIService} from '../../../APIService';
import {DataService} from '../../../DataService';

export default class Settings extends Component {
    static navigationOptions = {
        title: 'Settings',
        header: null,
    };

    render() {
        return (
            <LinearGradient colors={['#00af3a', '#005baf']} style={styles.background}>
                <View maxWidth={600} style={styles.topView}>
                    <StatusBar barStyle={'light-content'} />
                    <Text style={styles.title}>Settings</Text>
                    <ScrollView style={styles.scrollView}>
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
                        <Card
                            text={'Avatar'}
                            textColor={'#444444'}
                            iconName={'image'}
                            iconType={'feather'}
                            iconColor={'#444444'}
                            showArrow={true}
                            height={70}
                            onPress={() =>
                                this.props.navigation.navigate('Edit', {
                                    name: 'Avatar',
                                })
                            }
                        />
                        <Card
                            text={'Email'}
                            textColor={'#444444'}
                            iconName={'at-sign'}
                            iconType={'feather'}
                            iconColor={'#444444'}
                            showArrow={true}
                            height={70}
                            onPress={() =>
                                this.props.navigation.navigate('Edit', {
                                    name: 'Email Address',
                                })
                            }
                        />
                        <Card
                            text={'Password'}
                            textColor={'#444444'}
                            iconName={'key'}
                            iconType={'feather'}
                            iconColor={'#444444'}
                            showArrow={true}
                            height={70}
                            onPress={() =>
                                this.props.navigation.navigate('Edit', {
                                    name: 'Password',
                                })
                            }
                        />
                        <View style={styles.spacer} />
                        <Card
                            text={'Delete Account'}
                            textColor={'#ff4444'}
                            iconName={'trash'}
                            iconType={'feather'}
                            iconColor={'#ff4444'}
                            showArrow={false}
                            height={50}
                            onPress={() =>
                                DataService.getUsername().then(username =>
                                    APIService.deleteAccount(username)
                                        .then(res => res.json())
                                        .then(j => {
                                            console.log(j, j.status);
                                            if (j.status === 'success') {
                                                console.log('success');
                                                DataService.removeUserToken();
                                                this.props.navigation.navigate('Auth');
                                            } else {
                                                console.log('error');
                                            }
                                        }),
                                )
                            }
                        />
                        <Card
                            text={'Remove token'}
                            textColor={'#ff4444'}
                            iconName={'trash'}
                            iconType={'feather'}
                            iconColor={'#ff4444'}
                            showArrow={false}
                            height={50}
                            onPress={() => DataService.removeUserToken()}
                        />
                        <Card
                            text={'Sign Out'}
                            textColor={'#ff4444'}
                            iconName={'log-out'}
                            iconType={'feather'}
                            iconColor={'#ff4444'}
                            showArrow={false}
                            height={50}
                            onPress={() => this._signOutAsync()}
                        />
                        <View style={styles.versionContainer}>
                            <Text style={styles.versionText}>v0.0.1 alpha</Text>
                        </View>
                    </ScrollView>
                </View>
            </LinearGradient>
        );
    }

    async _signOutAsync() {
        if ((await APIService.logout()) === 'success') {
            await DataService.removeUserToken();
            this.props.navigation.navigate('Auth');
        } else {
            console.log('error');
        }
    }
}

export class Card extends Component {
    render() {
        return (
            <TouchableOpacity onPress={() => this.props.onPress()}>
                <View style={[styles.card, {height: this.props.height}]}>
                    <LinearGradient
                        start={{x: 0, y: 0}}
                        end={{x: 1, y: 0}}
                        colors={['#ffffff', '#dddddd']}
                        style={[styles.grad, {height: this.props.height}]}>
                        <View style={styles.cardIcon}>
                            <Icon name={this.props.iconName} type={this.props.iconType} color={this.props.iconColor} />
                        </View>
                        <View style={styles.cardText}>
                            <Text style={[styles.text, {color: this.props.textColor}]}>{this.props.text}</Text>
                        </View>
                        <View style={styles.cardArrow}>
                            {this.props.showArrow ? (
                                <Icon name={'chevron-right'} color={'#444444'} type={'feather'} />
                            ) : null}
                        </View>
                    </LinearGradient>
                </View>
            </TouchableOpacity>
        );
    }
}

export const styles = StyleSheet.create({
    title: {
        marginLeft: 20,
        marginTop: 0,
        marginBottom: 10,
        fontSize: 40,
        color: '#ffffff',
        fontWeight: '800',
    },
    topView: {
        position: 'relative',
        width: '100%',
        height: '108%',
        flex: 1,
    },
    scrollView: {
        position: 'relative',
        width: '100%',
        height: '91%',
    },
    background: {
        display: 'flex',
        flexDirection: 'row',
        justifyContent: 'center',
        paddingTop: '15%',
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
    text: {
        fontSize: 20,
        fontWeight: '600',
    },
    container: {
        position: 'relative',
        width: '100%',
        height: '100%',
    },
    spacer: {
        height: 20,
    },
    versionContainer: {
        justifyContent: 'center',
        alignItems: 'center',
        paddingTop: 10,
        paddingBottom: 10,
    },
    versionText: {
        color: '#bbbbbb',
        marginBottom: 10,
    },
    cardIcon: {
        position: 'absolute',
        left: 20,
    },
    cardText: {
        position: 'absolute',
        left: 60,
    },
    cardArrow: {
        position: 'absolute',
        right: 10,
    },
});
