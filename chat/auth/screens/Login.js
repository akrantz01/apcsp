import React, {Component} from 'react';
import {
    View,
    StyleSheet,
    Button,
    StatusBar,
    Text,
    Keyboard,
    TouchableWithoutFeedback,
} from 'react-native';
import AsyncStorage from '@react-native-community/async-storage';
import {Input} from 'react-native-elements';
import LinearGradient from 'react-native-linear-gradient';
import {Button as NButton} from 'native-base';
import {KeyboardAwareScrollView} from 'react-native-keyboard-aware-scroll-view';

export class Login extends Component {
    static navigationOptions = {
        title: 'Please sign in',
        header: null,
    };

    render() {
        return (
            <TouchableWithoutFeedback
                onPress={Keyboard.dismiss}
                accessible={false}>
                <View style={styles.view}>
                    <StatusBar barStyle={'light-content'} />
                    <LinearGradient
                        colors={['#0a49bf', '#182a4d']}
                        style={styles.background}>
                        <KeyboardAwareScrollView scrollEnabled={false}>
                            <Text style={styles.text}>Login</Text>
                            <View style={styles.card}>
                                <Input
                                    placeholder={'Username'}
                                    placeholderTextColor={'#BBB'}
                                    inputStyle={{color: 'white'}}
                                />
                                <View style={styles.spacer} />
                                <Input
                                    placeholder={'Password'}
                                    placeholderTextColor={'#BBB'}
                                    secureTextEntry={true}
                                    inputStyle={{color: 'white'}}
                                />
                                <NButton
                                    style={styles.button}
                                    title="Authenticate"
                                    onPress={this._signInAsync}>
                                    <LinearGradient
                                        start={{x: 0, y: 0}}
                                        end={{x: 1, y: 0}}
                                        colors={['#FF512F', '#F09819']}
                                        style={styles.buttonGrad}>
                                        <Text style={styles.buttonText}>
                                            Go
                                        </Text>
                                    </LinearGradient>
                                </NButton>
                            </View>
                            <Button
                                title="Don't have an account?"
                                color={'#FFFFFF'}
                                onPress={this._signUpAsync}
                            />
                        </KeyboardAwareScrollView>
                    </LinearGradient>
                </View>
            </TouchableWithoutFeedback>
        );
    }
    _signUpAsync = async () => {
        this.props.navigation.navigate('Signup');
    };

    _signInAsync = async () => {
        await AsyncStorage.setItem('userToken', 'abc');
        this.props.navigation.navigate('App');
    };
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
