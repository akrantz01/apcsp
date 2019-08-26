import React from 'react';
import {Component} from 'react';
import {
    RefreshControl,
    SectionList,
    StyleSheet,
    Text,
    View,
    TouchableOpacity,
    StatusBar,
    SafeAreaView,
} from 'react-native';
import LinearGradient from 'react-native-linear-gradient';
import {Avatar, Icon} from 'react-native-elements';

export default class MessagesScreen extends Component {
    constructor(props) {
        super(props);
        this.state = {
            refreshing: false,
        };
    }

    static navigationOptions = {
        title: 'Messages',
        header: null,
    };

    _onRefresh() {
        this.setState({refreshing: true});
        // .then( () => {
        this.setState({refreshing: false});
        // });
    }

    getData() {
        return [
            {data: [{key: 0, name: 'Alex Beaver', unread: true}]},
            {data: [{key: 1, name: 'Alex Krantz', unread: true}]},
            {data: [{key: 2, name: 'Aidan Sacco', unread: false}]},
            {data: [{key: 3, name: 'Guy Wilks', unread: true}]},
            {data: [{key: 4, name: 'Daniel Longo', unread: false}]},
            {data: [{key: 5, name: 'Eric Ettlin', unread: false}]},
            {data: [{key: 6, name: 'Dylan Pratt', unread: false}]},
        ];
    }

    render() {
        const {navigate} = this.props.navigation;
        return (
            <LinearGradient
                style={styles.background}
                colors={['#00af3a', '#005baf'] /*['#ee6210', '#cc0059']*/}>
                <StatusBar barStyle={'light-content'} />
                <Text style={styles.text}>Messages</Text>
                <View style={styles.buttonBar}>
                    <View
                        style={[
                            styles.buttonBarButtonContainer,
                            {marginLeft: 0, marginRight: 5, flex: 4},
                        ]}>
                        <TouchableOpacity onPress={null}>
                            <View style={styles.buttonBarButtonBackground}>
                                <View style={styles.buttonBarContents}>
                                    <Icon
                                        name={'plus'}
                                        type={'feather'}
                                        color={'#ffffff'}
                                    />
                                    <Text
                                        style={{
                                            marginLeft: 5,
                                            color: '#ffffff',
                                        }}>
                                        New Message
                                    </Text>
                                </View>
                            </View>
                        </TouchableOpacity>
                    </View>
                    <View
                        style={[
                            styles.buttonBarButtonContainer,
                            {marginLeft: 5, marginRight: 0, flex: 3},
                        ]}>
                        <TouchableOpacity onPress={() => navigate('Settings')}>
                            <View style={styles.buttonBarButtonBackground}>
                                <View style={styles.buttonBarContents}>
                                    <Icon
                                        name={'settings'}
                                        type={'feather'}
                                        color={'#ffffff'}
                                    />
                                    <Text
                                        style={{
                                            marginLeft: 5,
                                            color: '#ffffff',
                                        }}>
                                        Settings
                                    </Text>
                                </View>
                            </View>
                        </TouchableOpacity>
                    </View>
                </View>
                <SafeAreaView style={{height: '88%'}}>
                    <SectionList
                        refreshControl={
                            <RefreshControl
                                refreshing={this.state.refreshing}
                                onRefresh={this._onRefresh.bind(this)}
                                tintColor={'#ffffff'}
                            />
                        }
                        sections={this.getData()}
                        contentContainerStyle={{paddingBottom: 10}}
                        renderItem={({item}) => (
                            <TouchableOpacity
                                onPress={() =>
                                    navigate('Thread', {
                                        name: item.name.split(' ')[0],
                                    })
                                }>
                                <View style={styles.card}>
                                    <LinearGradient
                                        start={{x: 0, y: 0}}
                                        end={{x: 1, y: 0}}
                                        colors={['#ffffff', '#dddddd']}
                                        style={{
                                            padding: 20,
                                            borderRadius: 30,
                                            justifyContent: 'center',
                                        }}>
                                        <Avatar rounded title="MD" size={50} />
                                        <View style={styles.cardText}>
                                            <Text style={styles.name}>
                                                {item.name}
                                            </Text>
                                            <Text
                                                style={{
                                                    marginTop: 5,
                                                    color: '#666666',
                                                    fontSize: 16,
                                                }}
                                                numberOfLines={1}>
                                                this is the previous message
                                            </Text>
                                        </View>
                                    </LinearGradient>
                                </View>
                            </TouchableOpacity>
                        )}
                        getItem={null}
                        getItemCount={null}
                    />
                </SafeAreaView>
            </LinearGradient>
        );
    }
}

const styles = StyleSheet.create({
    background: {
        backgroundColor: '#EEEEEE',
        paddingTop: '15%',
        paddingBottom: 0,
        height: '100%',
    },
    buttonBar: {
        height: 45,
        padding: 0,
        marginLeft: 20,
        marginRight: 20,
        marginBottom: 10,
        flexDirection: 'row',
        justifyContent: 'center',
        alignItems: 'center',
        display: 'flex',
    },
    buttonBarButtonContainer: {
        height: '100%',
        borderRadius: 40,
    },
    buttonBarButtonBackground: {
        height: '100%',
        borderRadius: 40,
        backgroundColor: '#4616b7',
    },
    buttonBarContents: {
        padding: 10,
        flexDirection: 'row',
        justifyContent: 'center',
        alignItems: 'center',
    },
    avatar: {
        position: 'absolute',
        left: 0,
    },
    text: {
        marginLeft: 20,
        marginBottom: 10,
        fontSize: 40,
        color: '#ffffff',
        fontWeight: '800',
    },
    cardText: {
        padding: 20,
        position: 'absolute',
        right: 0,
        top: 0,
        bottom: 0,
        left: 70,
    },
    name: {
        fontSize: 20,
        fontWeight: '600',
    },
    card: {
        marginLeft: 20,
        marginRight: 20,
        marginBottom: 10,
        marginTop: 10,
        borderRadius: 30,
        backgroundColor: '#FFFFFF',
        shadowColor: '#444444',
        shadowOffset: {width: 3, height: 6},
        shadowOpacity: 0.5,
        shadowRadius: 3,
        elevation: 1,
    },
});
