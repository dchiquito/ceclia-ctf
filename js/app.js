
app = new Vue({
    el: '.app',
    methods: {
        signout: function(e) {
            console.log('Signing out');
            $cookies.remove('auth');
            $cookies.remove('username');
            $cookies.remove('password');
        },
        hint: function(e) {
            console.log("Requesting hint");
        },
        submit: function(e) {
            console.log("Submitting a flag");
        }
    },
});


