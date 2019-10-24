import React, {Component} from 'react';
import {View, StyleSheet, Button, StatusBar, Text, Keyboard, TouchableWithoutFeedback} from 'react-native';
import {Input} from 'react-native-elements';
import LinearGradient from 'react-native-linear-gradient';
import {Button as NButton} from 'native-base';
import {KeyboardAwareScrollView} from 'react-native-keyboard-aware-scroll-view';
import {sha256} from 'react-native-sha256';
import {APIService} from '../../APIService';
import {DataService} from '../../DataService';
import Dialog from 'react-native-dialog';

export default class Login extends Component {
    constructor(props) {
        super(props);
        this.state = {
            dialogVisible: false,
            buttonEnabled: false,
            loading: false,
        };
    }

    static navigationOptions = {
        title: 'Please sign in',
        header: null,
    };

    render() {
        return (
            <TouchableWithoutFeedback onPress={Keyboard.dismiss} accessible={false}>
                <View style={styles.view}>
                    <Dialog.Container visible={this.state.dialogVisible}>
                        <Dialog.Title>Error</Dialog.Title>
                        <Dialog.Description>
                            Could not communicate with the server. Please try again later.
                        </Dialog.Description>
                        <Dialog.Button label="OK" bold onPress={() => this.setState({dialogVisible: false})} />
                    </Dialog.Container>

                    <StatusBar barStyle={'light-content'} />
                    <LinearGradient colors={['#0a49bf', '#182a4d']} style={styles.background}>
                        <KeyboardAwareScrollView scrollEnabled={false}>
                            <Text style={styles.text}>Login</Text>
                            <View style={styles.card}>
                                <Input
                                    placeholder={'Username'}
                                    placeholderTextColor={'#BBB'}
                                    inputStyle={styles.insideText}
                                    onChangeText={value => {
                                        this.setState({username: value});
                                        this.updateButton();
                                    }}
                                />
                                <View style={styles.spacer} />
                                <Input
                                    placeholder={'Password'}
                                    placeholderTextColor={'#BBB'}
                                    secureTextEntry={true}
                                    inputStyle={styles.insideText}
                                    onChangeText={value => {
                                        this.setState({password: value});
                                        this.updateButton();
                                    }}
                                />
                                <NButton
                                    style={styles.button}
                                    title="Go"
                                    disabled={!this.state.buttonEnabled || this.state.loading}
                                    onPress={() => this._signInAsync()}>
                                    <LinearGradient
                                        start={{x: 0, y: 0}}
                                        end={{x: 1, y: 0}}
                                        colors={['#FF512F', '#F09819']}
                                        style={styles.buttonGrad}>
                                        <Text style={styles.buttonText}>Go</Text>
                                    </LinearGradient>
                                </NButton>
                            </View>
                            <Button
                                title="Don't have an account?"
                                color={'#FFFFFF'}
                                disabled={this.state.loading}
                                onPress={() => this.props.navigation.navigate('Signup')}
                            />
                        </KeyboardAwareScrollView>
                    </LinearGradient>
                </View>
            </TouchableWithoutFeedback>
        );
    }

    updateButton() {
        if (!!this.state.password && !!this.state.username) {
            this.setState({buttonEnabled: true});
        } else {
            this.setState({buttonEnabled: false});
        }
    }

    async _signInAsync() {
        this.setState({loading: true});
        sha256(this.state.password).then(hash => {
            console.log('hash', hash);
            APIService.login(this.state.username, hash).then(succeeded => {
                if (succeeded) {
                    DataService.saveUsername(this.state.username);
                    this.props.navigation.navigate('App');
                } else {
                    console.log('error');
                    console.log(this);
                    this.setState({dialogVisible: true, loading: false, buttonEnabled: true});
                }
            });
        });
        // AsyncStorage.setItem('authToken', 'abc');
        // this.props.navigation.navigate('App');
    }
}

const styles = StyleSheet.create({
    background: {
        backgroundColor: '#EEEEEE',
        paddingTop: '40%',
        height: '100%',
        display: 'flex',
    },
    spacer: {
        marginTop: 10,
        marginBottom: 10,
    },
    text: {
        paddingLeft: 30,
        fontSize: 40,
        color: '#FFFFFF',
        fontWeight: '800',
    },
    insideText: {
        color: 'white',
    },
    card: {
        borderRadius: 12,
        padding: 20,
        margin: 20,
    },
    button: {
        color: '#000000',
        alignContent: 'center',
        borderRadius: 40,
        marginLeft: 10,
        marginRight: 10,
        marginTop: 40,
        padding: 0,
    },
    buttonGrad: {
        margin: 0,
        alignSelf: 'center',
        position: 'absolute',
        left: 0,
        top: 0,
        right: 0,
        bottom: 0,
        borderRadius: 40,
        justifyContent: 'center',
        alignItems: 'center',
    },
    buttonText: {
        fontSize: 18,
        color: 'white',
        fontWeight: 'bold',
    },
});
