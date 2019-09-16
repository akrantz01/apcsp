import React, {Component} from 'react';
import {View, StyleSheet, StatusBar, Text, TouchableWithoutFeedback, TouchableOpacity, Keyboard} from 'react-native';
import {Input, Icon} from 'react-native-elements';
import LinearGradient from 'react-native-linear-gradient';
import {KeyboardAwareScrollView} from 'react-native-keyboard-aware-scroll-view';
import {sha256} from 'react-native-sha256';

export class Signup extends Component {
    constructor(props) {
        super(props);
        this.state = {
            loading: false,
            fullName: '',
            email: '',
            username: '',
            password: '',
            confirmPass: '',
            keyboardOpen: false,
        };
        this._keyboardDidShow = this._keyboardDidShow.bind(this);
        this._keyboardDidHide = this._keyboardDidHide.bind(this);
    }

    componentDidMount() {
        this.keyboardDidShowListener = Keyboard.addListener('keyboardDidShow', this._keyboardDidShow);
        this.keyboardDidHideListener = Keyboard.addListener('keyboardDidHide', this._keyboardDidHide);
    }

    componentWillUnmount() {
        this.keyboardDidShowListener.remove();
        this.keyboardDidHideListener.remove();
    }

    _keyboardDidShow() {
        this.setState({keyboardOpen: true});
    }

    _keyboardDidHide() {
        this.setState({keyboardOpen: false});
    }

    static navigationOptions = {
        title: 'Please sign in',
        header: null,
    };

    async _signUpAsync() {
        console.log('$$$$$$$$$$$$$$$$$$$$');
        console.log(
            this.state.fullName,
            this.state.email,
            this.state.username,
            this.state.password,
            this.state.confirmPass,
        );
        console.log('start');
        this.setState({loading: true});
        if (this.state.password === this.state.confirmPass) {
            sha256(this.state.password).then(hash => {
                console.log('hash', hash);
                // APIService.register(
                //     this.state.fullName,
                //     this.state.email,
                //     this.state.username,
                //     hash,
                // ).then(res => {
                //     console.log('response', res);
                //     this.setState({loading: false});
                //     this.props.navigation.navigate('App');
                // });
            });
        }
    }

    render() {
        return (
            <TouchableWithoutFeedback onPress={Keyboard.dismiss} accessible={false}>
                <React.Fragment>
                    <StatusBar barStyle={'light-content'} />
                    <LinearGradient colors={['#460bbd', '#2a184d']} style={styles.background}>
                        <View style={styles.container}>
                            <View style={styles.back}>
                                <TouchableOpacity onPress={() => this.props.navigation.goBack()}>
                                    <Icon raised name="arrow-left" type="feather" color="#f50" />
                                </TouchableOpacity>
                            </View>
                            <Text style={styles.registerText}>Register</Text>
                            <KeyboardAwareScrollView
                                scrollEnabled={this.state.keyboardOpen}
                                style={styles.scrollView}
                                contentContainerStyle={styles.innerScrollView}>
                                <View style={styles.registerContainer}>
                                    <View style={styles.textBoxContainer}>
                                        <Input
                                            placeholder={'Full Name'}
                                            placeholderTextColor={'#CCCCCC'}
                                            inputStyle={styles.inputText}
                                            onChangeText={value => this.setState({fullName: value})}
                                        />
                                    </View>
                                    <View style={styles.textBoxContainer}>
                                        <Input
                                            placeholder={'E-mail'}
                                            placeholderTextColor={'#CCCCCC'}
                                            inputStyle={styles.inputText}
                                            onChangeText={value => this.setState({email: value})}
                                        />
                                    </View>
                                    <View style={styles.textBoxContainer}>
                                        <Input
                                            placeholder={'Username'}
                                            placeholderTextColor={'#CCCCCC'}
                                            inputStyle={styles.inputText}
                                            onChangeText={value => this.setState({username: value})}
                                        />
                                    </View>
                                    <View style={styles.textBoxContainer}>
                                        <Input
                                            placeholder={'Password'}
                                            placeholderTextColor={'#CCCCCC'}
                                            secureTextEntry={true}
                                            inputStyle={styles.inputText}
                                            onChangeText={value => this.setState({password: value})}
                                        />
                                    </View>
                                    <View style={styles.textBoxContainer}>
                                        <Input
                                            placeholder={'Confirm Password'}
                                            placeholderTextColor={'#CCCCCC'}
                                            secureTextEntry={true}
                                            inputStyle={styles.inputText}
                                            onChangeText={value =>
                                                this.setState({
                                                    confirmPass: value,
                                                })
                                            }
                                        />
                                    </View>
                                    <TouchableOpacity
                                        style={styles.goButton}
                                        disabled={this.state.loading}
                                        onPress={() => this._signUpAsync()}>
                                        <LinearGradient
                                            start={{x: 0, y: 0}}
                                            end={{x: 1, y: 0}}
                                            colors={['#FF512F', '#F09819']}
                                            style={styles.buttonGrad}>
                                            <Text style={styles.buttonText}>Sign Up</Text>
                                        </LinearGradient>
                                    </TouchableOpacity>
                                </View>
                            </KeyboardAwareScrollView>
                        </View>
                    </LinearGradient>
                </React.Fragment>
            </TouchableWithoutFeedback>
        );
    }
}

const styles = StyleSheet.create({
    background: {
        width: '100%',
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        alignContent: 'flex-end',
        justifyContent: 'flex-end',
    },
    container: {
        height: '95%',
        width: '100%',
        alignContent: 'center',
        justifyContent: 'center',
    },
    back: {
        width: 0,
        margin: '3%',
    },
    registerText: {
        fontSize: 40,
        color: '#FFFFFF',
        fontWeight: '800',
        width: '90%',
        alignSelf: 'center',
    },
    scrollView: {
        width: '100%',
        height: '100%',
    },
    innerScrollView: {
        display: 'flex',
        flexDirection: 'column',
        alignContent: 'center',
        justifyContent: 'center',
    },
    registerContainer: {
        width: '80%',
        alignSelf: 'center',
        justifyContent: 'center',
    },
    inputText: {
        color: 'white',
    },
    textBoxContainer: {
        marginTop: 10,
        marginBottom: 0,
    },
    goButton: {
        color: '#000000',
        marginTop: 25,
        borderRadius: 40,
        width: '100%',
        height: 50,
        alignSelf: 'center',
        justifyContent: 'center',
        alignContent: 'center',
    },
    buttonGrad: {
        width: '100%',
        height: '100%',
        borderRadius: 40,
        justifyContent: 'center',
        alignContent: 'center',
    },
    buttonText: {
        fontSize: 18,
        color: 'white',
        fontWeight: 'bold',
        alignSelf: 'center',
    },
});
