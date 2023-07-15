package org.bjason.bettermaths;

import androidx.appcompat.app.AppCompatActivity;

import android.os.Bundle;

import org.bjason.bettermaths.mobilegame.EbitenView;

import go.Seq;


public class MainActivity extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.ebiten);
        Seq.setContext(getApplicationContext());

    }
    private EbitenView getEbitenView() {
        return (EbitenView)this.findViewById(R.id.mobile_game);
    }

    @Override
    protected void onPause() {
        super.onPause();
        this.getEbitenView().suspendGame();
    }

    @Override
    protected void onResume() {
        super.onResume();
        EbitenView fred = this.getEbitenView();
        fred.resumeGame();
    }

}