import React, {Component} from 'react';
import {
    View,
    StyleSheet,
    StatusBar,
    Text,
    TouchableWithoutFeedback,
    TouchableOpacity,
    Keyboard,
} from 'react-native';
import {Input, Icon} from 'react-native-elements';
import LinearGradient from 'react-native-linear-gradient';
import {Button} from 'native-base';
import {KeyboardAwareScrollView} from 'react-native-keyboard-aware-scroll-view';

export class Signup extends Component {
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
                        colors={['#460bbd', '#2a184d']}
                        style={styles.background}>
                        <KeyboardAwareScrollView scrollEnabled={false}>
                            <View style={styles.back}>
                                <TouchableOpacity
                                    onPress={() =>
                                        this.props.navigation.goBack()
                                    }>
                                    <Icon
                                        raised
                                        name="arrow-left"
                                        type="feather"
                                        color="#f50"
                                    />
                                </TouchableOpacity>
                            </View>
                            <Text style={styles.text}>Register</Text>
                            <View style={styles.card} behavior="margin" enabled>
                                <Input
                                    placeholder={'E-mail'}
                                    placeholderTextColor={'#BBBBBB'}
                                    style={styles.text}
                                    inputStyle={{color: 'white'}}
                                />
                                <View style={styles.spacer} />
                                <Input
                                    placeholder={'Username'}
                                    placeholderTextColor={'#CCCCCC'}
                                    style={styles.text}
                                    inputStyle={{color: 'white'}}
                                />
                                <View style={styles.spacer} />
                                <Input
                                    placeholder={'Password'}
                                    placeholderTextColor={'#BBBBBB'}
                                    style={styles.text}
                                    secureTextEntry={true}
                                    inputStyle={{color: 'white'}}
                                />
                                <View style={styles.spacer} />
                                <Input
                                    placeholder={'Confirm Password'}
                                    placeholderTextColor={'#BBBBBB'}
                                    style={styles.text}
                                    secureTextEntry={true}
                                    inputStyle={{color: 'white'}}
                                />
                                <Button
                                    style={styles.button}
                                    onPress={this._signUpAsync}>
                                    <LinearGradient
                                        start={{x: 0, y: 0}}
                                        end={{x: 1, y: 0}}
                                        colors={['#FF512F', '#F09819']}
                                        style={styles.buttonGrad}>
                                        <Text style={styles.buttonText}>
                                            Sign Up
                                        </Text>
                                    </LinearGradient>
                                </Button>
                            </View>
                        </KeyboardAwareScrollView>
                    </LinearGradient>
                </View>
            </TouchableWithoutFeedback>
        );
    }

    _signUpAsync = async () => {};
}

const styles = StyleSheet.create({
    background: {
        backgroundColor: '#EEEEEE',
        height: '100%',
        display: 'flex',
    },
    back: {
        marginTop: 40,
        marginLeft: 15,
        marginBottom: 35,
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
