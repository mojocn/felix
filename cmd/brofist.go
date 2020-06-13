package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/libragen/felix/util"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

// broFistCmd represents the brofist command
var broFistCmd = &cobra.Command{
	Use:   "brofist",
	Short: "Pewdiepie needs your help.Do your part to subscribe Pewdiepie's Youtube Channel.",
	Long:  `PewDiePie vs T-Series`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(strings.Replace(brofistChars, idchar, "`", -1))
		color.Red("Pewdiepie needs your help.")
		color.Yellow("Do your part to subscribe Felix's Youtube Channel.")
		color.Cyan("Your browser will go to Pewdiepie's Channel.")
		color.Blue("Please click the subscribe button on the right.")
		time.Sleep(time.Second * 2)
		util.BrowserOpen("https://www.youtube.com/channel/UC-lHJZR3Gqxm24_Vd_AJ5Yw")

	},
}

func init() {
	rootCmd.AddCommand(broFistCmd)
}

const idchar = "[eTe]"
const brofistChars = `
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMmhs+////oyhmMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMNh/.-/oyyys+/-.:ohmMMMNdyo//+oydNMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMNh:[eTe]:hNMNdddmNMNds/-.:/---:+sys+:--odMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMmdhssssssyhy:[eTe]:hNMm+.[eTe][eTe][eTe][eTe]./sdNMNdsyhdmNNddmNNms-[eTe]+mMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMNy:.-:/++++//::+hNMmo.         [eTe]./sdddhyo/-[eTe][eTe].-odMNs..yMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMh-[eTe]+dNNmmmmmmNNMMNms.               [eTe][eTe][eTe]          [eTe]oNMm:[eTe]oMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMNdso/////[eTe]:dMNs-......-:++:[eTe]                                sMMN: yMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMNs-.:osyysoyNNh-                                              .NMMd .NMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMm:[eTe]omMNmddmmmh/                                                 hMMM+ oMMMMMMMMMMMMMMMMM
MMMMMMMMMMMN: sMMh-[eTe][eTe][eTe][eTe]..                                                   /MMMN. hMMMMMMMMMMMMMMMM
MMMMMMMMMMMy /MMd[eTe]                                                          .MMMMd[eTe].NMMMMMMMMMMMMMMM
MMMMMMMMMMM. dMM-                                          /y/               NMMMMo +MMMMMMMMMMMMMMM
MMMMMMMMMMd -MMm       [eTe][eTe]                                  hMN               hMMmMN. dMMMMMMMMMMMMMM
MMMMMMMMMMo oMMo      :dd-                                 sMM-              oMMsNMh -NMMMMMMMMMMMMM
MMMMMMMMMM- dMM:      /MMs             -+-                 :MMo              /MMssMM/ sMMMMMMMMMMMMM
MMMMMMMMMN [eTe]NMN[eTe]      -MMh             NMm                 [eTe]MMh              -MMh[eTe]NMN[eTe][eTe]mMMMMMMMMMMMM
MMMMMMMMMh :MMh       -MMd             NMm                  mMN[eTe]             [eTe]MMm +MMo +MMMMMMMMMMMM
MMMMMMMMMo oMMo       -MMh            [eTe]NMN                  sMM:             [eTe]NMN [eTe]NMN[eTe][eTe]NMMMMMMMMMMM
MMMMMMMMM/ yMM/       /MMy            [eTe]NMN                  /MMy              mMM[eTe] sMM+ sMMMMMMMMMMM
MMMMMMMMM/ yMM/       oMMo            [eTe]MMN                  [eTe]NMN              hMM- .MMd -MMMMMMMMMMM
MMMMMMMMM+ sMM+       hMM:            .MMm                   mMM.             sMM+  dMM- dMMMMMMMMMM
MMMMMMMMMs +MMs       mMM[eTe]            -MMd                   dMM-             +MMs  /MMs +MMMMMMMMMM
MMMMMMMMMm .NMN:     [eTe]NMN             -MMd                   dMM-             :MMh  [eTe]NMN[eTe][eTe]NMMMMMMMMM
MMMMMMMMMMo :mMNh/.  -MMd             :MMh                   dMM-             :MMh   sMM/ yMMMMMMMMM
MMMMMMMMMMNy.[eTe]+dNNNdsyMMh             :MMh                   mMM.             yMM+   :MMs +MMMMMMMMM
MMMMMMMMMMMMNy:..+ydmNMMN.            -MMd                  [eTe]NMN            [eTe]sNMh[eTe]   sMM+ sMMMMMMMMM
MMMMMMMMMMMMMMNmy+:--./NMmo:-[eTe]        .NMN.                 .MMN.[eTe]    [eTe][eTe].-/smMNs[eTe]  [eTe]yNMy[eTe].NMMMMMMMMM
MMMMMMMMMMMMMMMMMMNNmo[eTe]-sdNNNmdhyyyyhdmMMMm+-[eTe]          [eTe].-:hMMMmmmddmmmmNNNds-   .dMNo[eTe]:mMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMdo---:oyhdddddhhs+/ymMNmdhyssssyhdmmNNNdo+syyhyyyss+/-[eTe]    :mMm/ +NMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMNmdyo/::::::::/o+-.:oyhddmMMMmhyso+:-[eTe]                  /NMd-[eTe]sNMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMNNNmNNNMMMMmds+/:.[eTe]/MMd.                         +NMh.[eTe]hMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMNs +MMs                         +NMh[eTe].dMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMd -MMd[eTe]                       oNMh[eTe].dMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM: sMMy-[eTe][eTe][eTe]..:/+ooo++//:-....oMMy[eTe].dMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMm-[eTe]omMNddmmNNNmmmmmmNNNNNNNNNms[eTe]-mMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMNs../syyso/::-------:::/++/:-.+mMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMNhs+//+osyhdmmmmdddhyysssyhNMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM
`
