import React, {Component} from 'react';
import {GiftedChat} from 'react-native-gifted-chat';
import LinearGradient from 'react-native-linear-gradient';
import {StyleSheet, View, TouchableOpacity, SafeAreaView, Text} from 'react-native';
import {Input, Icon} from 'react-native-elements';
import {DataService} from '../../../DataService';

export default class Thread extends Component {
    static navigationOptions = ({navigation}) => {
        return {
            title: navigation.getParam('name', 'null'),
            header: null,
        };
    };

    state = {
        messages: DataService.getSavedMessages(this.props.navigation.getParam('chatID', 'null')),
        compose: '',
    };

    onSend(messages) {
        if (this.state.compose !== '') {
            this.setState(previousState => ({
                messages: GiftedChat.append(previousState.messages, {
                    text: messages.text,
                    _id: previousState.messages[0]._id + 1,
                    createdAt: new Date(),
                    user: {
                        _id: 1,
                        name: 'React Native',
                        avatar: 'https://placeimg.com/140/140/any',
                    },
                }),
            }));
            this.setState({compose: ''});
            DataService.saveMessage(this.props.navigation.getParam('chatID'), this.state.messages);
        }
    }

    renderInputToolbar() {
        return (
            <View>
                <View style={styles.composeBar}>
                    <View>
                        <Input
                            multiline={true}
                            value={this.state.compose}
                            onChangeText={value => this.setState({compose: value})}
                        />
                    </View>
                    <TouchableOpacity
                        style={styles.button}
                        disabled={!this.state.compose}
                        onPress={() => this.onSend({text: this.state.compose.trim()})}>
                        <Icon
                            reverse
                            reverseColor={'#ffffff'}
                            size={20}
                            type={'feather'}
                            name={'send'}
                            color={this.state.compose !== '' ? '#0083ff' : '#838383'}
                        />
                    </TouchableOpacity>
                </View>
            </View>
        );
    }

    render() {
        return (
            <LinearGradient style={styles.background} colors={['#00af3a', '#005baf']}>
                <SafeAreaView style={styles.safeArea} maxWidth={600}>
                    <View style={styles.topBar}>
                        <View style={styles.backButton}>
                            <Icon
                                raised
                                size={20}
                                name={'arrow-left'}
                                type={'feather'}
                                color={'#444444'}
                                onPress={() =>
                                    this.props.navigation.dispatch({
                                        type: 'Navigation/BACK',
                                    })
                                }
                            />
                        </View>
                        <View style={styles.textBox}>
                            <Text style={styles.text} numberOfLines={1}>
                                {this.props.navigation.getParam('name', 'Message')}
                            </Text>
                        </View>
                    </View>
                    <GiftedChat
                        alignTop={false}
                        renderInputToolbar={() => this.renderInputToolbar()}
                        messages={this.state.messages}
                        onSend={messages => this.onSend(messages)}
                        user={{
                            _id: 1,
                        }}
                    />
                </SafeAreaView>
            </LinearGradient>
        );
    }
}

const styles = StyleSheet.create({
    background: {
        backgroundColor: '#EEEEEE',
        height: '100%',
        width: '100%',
        display: 'flex',
        flexDirection: 'row',
        paddingBottom: 25,
        justifyContent: 'center',
    },
    safeArea: {
        flex: 1,
        display: 'flex',
    },
    backButton: {
        marginLeft: -5,
    },
    topBar: {
        display: 'flex',
        flexDirection: 'row',
        marginLeft: 20,
        marginRight: 20,
        borderBottomColor: '#CCC',
        paddingBottom: 5,
        borderBottomWidth: 2,
    },
    composeBar: {
        width: '90%',
        alignSelf: 'center',
        paddingRight: 45,
        paddingLeft: 10,
        paddingTop: 5,
        marginTop: 5,
        height: 50,

        borderRadius: 30,
        backgroundColor: '#FFFFFF',
        shadowColor: '#444444',
        shadowOffset: {width: 3, height: 6},
        shadowOpacity: 0.5,
        shadowRadius: 3,
        elevation: 1,
    },
    button: {
        position: 'absolute',
        right: -4,
        bottom: -4,
    },
    textBox: {
        display: 'flex',
        flexDirection: 'column',
        alignContent: 'flex-end',
        justifyContent: 'flex-end',
        flex: 1,
    },
    text: {
        fontSize: 30,
        alignSelf: 'flex-end',
        color: 'white',
        fontWeight: '800',
        marginRight: 5,
    },
});
