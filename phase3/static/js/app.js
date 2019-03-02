
app = new Vue({
    el: '.app',
    data: {
        flag: ''
    },
    methods: {
        signout: function(e) {
            console.log('Signing out');
            $cookies.remove('auth');
        },
        reset: function() {
            console.log('Reseting all progress');
            location.href = '/app/reset';
        }
    }
});


