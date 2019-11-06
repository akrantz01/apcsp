import React from 'react';
import {Component} from 'react';
import {View, StatusBar, Text, StyleSheet} from 'react-native';
import LinearGradient from 'react-native-linear-gradient';
import {Card, styles} from './Settings';
import {Input, Icon} from 'react-native-elements';
import {DataService} from '../../../DataService';
export default class Edit extends Component {
    constructor(props) {
        super(props);
        this.state = {
            newVal: '',
            confPass: '',
        };
    }
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
                    <View style={[styles.card, {height: this.props.height}]}>
                        <LinearGradient
                            start={{x: 0, y: 0}}
                            end={{x: 1, y: 0}}
                            colors={['#ffffff', '#dddddd']}
                            style={[styles.grad, {height: 70}]}>
                            <View style={styles.cardIcon}>
                                <Icon
                                    name={this.props.navigation.getParam('icon', '')}
                                    type={'feather'}
                                    color={'#444444'}
                                />
                            </View>
                            <View style={[styles.cardText, thisStyles.textInput]}>
                                <Input
                                    placeholder={'New ' + this.props.navigation.getParam('name', 'null')}
                                    placeholderTextColor={'#BBB'}
                                    inputStyle={styles.insideText}
                                    onChangeText={value => {
                                        this.setState({newVal: value});
                                        this.updateButton();
                                    }}
                                />
                            </View>
                        </LinearGradient>
                    </View>
                    <View style={[styles.card, {height: this.props.height}]}>
                        <LinearGradient
                            start={{x: 0, y: 0}}
                            end={{x: 1, y: 0}}
                            colors={['#ffffff', '#dddddd']}
                            style={[styles.grad, {height: 70}]}>
                            <View style={styles.cardIcon}>
                                <Icon name={'key'} type={'feather'} color={'#444444'} />
                            </View>
                            <View style={[styles.cardText, thisStyles.textInput]}>
                                <Input
                                    placeholder={'Confirm Password'}
                                    placeholderTextColor={'#BBB'}
                                    inputStyle={styles.insideText}
                                    onChangeText={value => {
                                        this.setState({confPass: value});
                                        this.updateButton();
                                    }}
                                />
                            </View>
                        </LinearGradient>
                    </View>
                    <Card
                        text={'Confirm'}
                        textStyle={{fontSize: 25, textAlign: 'center'}}
                        textColor={'#35c454'}
                        iconName={''}
                        iconType={'feather'}
                        iconColor={'#35c454'}
                        showArrow={false}
                        height={60}
                        onPress={() => DataService.removeUserToken()}
                    />
                </View>
            </LinearGradient>
        );
    }
}

const thisStyles = StyleSheet.create({
    title: {
        paddingTop: 55,
    },
    textInput: {
        paddingTop: 20,
        paddingBottom: 20,
        width: '80%',
    },
});
